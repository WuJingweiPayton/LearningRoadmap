// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {IERC721} from "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import {IERC721Receiver} from "@openzeppelin/contracts/token/ERC721/IERC721Receiver.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC20Metadata} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";
import {AggregatorV3Interface} from "./interfaces/AggregatorV3Interface.sol";

/// @title Shared NFT auction logic for UUPS and transparent proxy implementations.
abstract contract AuctionHouseCore is Initializable, OwnableUpgradeable, ReentrancyGuardUpgradeable, IERC721Receiver {
    using SafeERC20 for IERC20;

    uint256 public constant BPS_DENOMINATOR = 10_000;
    uint256 public constant USD_DECIMALS = 18;

    struct Auction {
        address seller;
        address nft;
        uint256 tokenId;
        uint64 endTime;
        address highestBidder;
        uint256 highestBidAmount;
        address bidAsset;
        uint256 highestBidUsd;
        bool settled;
    }

    uint256 public auctionCount;
    address public feeRecipient;
    mapping(uint256 => Auction) public auctions;
    mapping(address => address) public priceFeeds;

    event AuctionCreated(
        uint256 indexed auctionId, address indexed seller, address indexed nft, uint256 tokenId, uint64 endTime
    );
    event BidPlaced(
        uint256 indexed auctionId, address indexed bidder, address asset, uint256 amount, uint256 usdAmount
    );
    event AuctionSettled(uint256 indexed auctionId, address indexed winner, address asset, uint256 amount, uint256 fee);
    event PriceFeedUpdated(address indexed asset, address indexed feed);

    error ZeroAddress();
    error InvalidDuration();
    error InvalidBid();
    error AuctionNotFound(uint256 auctionId);
    error AuctionAlreadySettled(uint256 auctionId);
    error AuctionEnded(uint256 auctionId);
    error AuctionOngoing(uint256 auctionId);
    error BidTooLow(uint256 bidUsd, uint256 currentUsd);
    error MissingPriceFeed(address asset);
    error BadPrice(address feed);
    error TransferFailed();

    function implementationVersion() public pure virtual returns (uint256) {
        return 1;
    }

    function setFeeRecipient(address newRecipient) external onlyOwner {
        if (newRecipient == address(0)) revert ZeroAddress();
        feeRecipient = newRecipient;
    }

    function setPriceFeed(address asset, address feed) external onlyOwner {
        if (feed == address(0)) revert ZeroAddress();
        priceFeeds[asset] = feed;
        emit PriceFeedUpdated(asset, feed);
    }

    function createAuction(address nft, uint256 tokenId, uint256 duration)
        external
        nonReentrant
        returns (uint256 auctionId)
    {
        if (nft == address(0)) revert ZeroAddress();
        if (duration == 0) revert InvalidDuration();
        uint256 end = block.timestamp + duration;
        if (end > type(uint64).max) revert InvalidDuration();

        auctionId = ++auctionCount;
        uint64 endTime = uint64(end);
        auctions[auctionId] = Auction({
            seller: msg.sender,
            nft: nft,
            tokenId: tokenId,
            endTime: endTime,
            highestBidder: address(0),
            highestBidAmount: 0,
            bidAsset: address(0),
            highestBidUsd: 0,
            settled: false
        });

        IERC721(nft).safeTransferFrom(msg.sender, address(this), tokenId);
        emit AuctionCreated(auctionId, msg.sender, nft, tokenId, endTime);
    }

    function bidWithEth(uint256 auctionId) external payable nonReentrant {
        _bid(auctionId, address(0), msg.value);
    }

    function bidWithErc20(uint256 auctionId, address token, uint256 amount) external payable nonReentrant {
        if (token == address(0) || msg.value != 0) revert InvalidBid();
        _bid(auctionId, token, amount);
    }

    function settle(uint256 auctionId) external nonReentrant {
        Auction storage auction = auctions[auctionId];
        _requireOpen(auctionId, auction);
        if (block.timestamp < auction.endTime) revert AuctionOngoing(auctionId);
        auction.settled = true;

        if (auction.highestBidder == address(0)) {
            IERC721(auction.nft).safeTransferFrom(address(this), auction.seller, auction.tokenId);
            emit AuctionSettled(auctionId, address(0), address(0), 0, 0);
            return;
        }

        IERC721(auction.nft).safeTransferFrom(address(this), auction.highestBidder, auction.tokenId);

        uint256 fee = (auction.highestBidAmount * feeBps(auction.highestBidUsd)) / BPS_DENOMINATOR;
        _payout(auction.bidAsset, feeRecipient, fee);
        _payout(auction.bidAsset, auction.seller, auction.highestBidAmount - fee);
        emit AuctionSettled(auctionId, auction.highestBidder, auction.bidAsset, auction.highestBidAmount, fee);
    }

    function quoteUsd(address asset, uint256 amount) public view returns (uint256) {
        address feed = priceFeeds[asset];
        if (feed == address(0)) revert MissingPriceFeed(asset);

        (uint80 roundId, int256 answer,, uint256 updatedAt, uint80 answeredInRound) =
            AggregatorV3Interface(feed).latestRoundData();
        if (answer <= 0 || updatedAt == 0 || answeredInRound < roundId) revert BadPrice(feed);

        uint8 priceDecimals = AggregatorV3Interface(feed).decimals();
        if (priceDecimals > USD_DECIMALS) revert BadPrice(feed);

        uint8 assetDecimals = asset == address(0) ? 18 : IERC20Metadata(asset).decimals();
        uint256 scaleUp = 10 ** (USD_DECIMALS - priceDecimals);
        uint256 scaleDown = 10 ** assetDecimals;
        return Math.mulDiv(amount, uint256(answer) * scaleUp, scaleDown);
    }

    function feeBps(uint256 usdAmount) public pure returns (uint256) {
        if (usdAmount < 100 * 10 ** USD_DECIMALS) return 250;
        if (usdAmount < 1000 * 10 ** USD_DECIMALS) return 150;
        return 50;
    }

    function onERC721Received(address, address, uint256, bytes calldata) external pure returns (bytes4) {
        return IERC721Receiver.onERC721Received.selector;
    }

    function _initializeAuctionHouse(address initialOwner, address ethUsdFeed, address feeRecipient_) internal {
        if (initialOwner == address(0) || ethUsdFeed == address(0) || feeRecipient_ == address(0)) {
            revert ZeroAddress();
        }
        __Ownable_init(initialOwner);
        __ReentrancyGuard_init();
        priceFeeds[address(0)] = ethUsdFeed;
        feeRecipient = feeRecipient_;
        emit PriceFeedUpdated(address(0), ethUsdFeed);
    }

    function _bid(uint256 auctionId, address asset, uint256 amount) internal {
        Auction storage auction = auctions[auctionId];
        _requireOpen(auctionId, auction);
        if (block.timestamp >= auction.endTime) revert AuctionEnded(auctionId);
        if (amount == 0) revert InvalidBid();
        if (asset == address(0) && msg.value != amount) revert InvalidBid();

        uint256 usdAmount = quoteUsd(asset, amount);
        if (usdAmount <= auction.highestBidUsd) revert BidTooLow(usdAmount, auction.highestBidUsd);

        address previousBidder = auction.highestBidder;
        address previousAsset = auction.bidAsset;
        uint256 previousAmount = auction.highestBidAmount;

        auction.highestBidder = msg.sender;
        auction.bidAsset = asset;
        auction.highestBidAmount = amount;
        auction.highestBidUsd = usdAmount;

        if (previousBidder != address(0)) {
            _payout(previousAsset, previousBidder, previousAmount);
        }
        if (asset != address(0)) {
            IERC20(asset).safeTransferFrom(msg.sender, address(this), amount);
        }

        emit BidPlaced(auctionId, msg.sender, asset, amount, usdAmount);
    }

    function _requireOpen(uint256 auctionId, Auction storage auction) internal view {
        if (auction.seller == address(0)) revert AuctionNotFound(auctionId);
        if (auction.settled) revert AuctionAlreadySettled(auctionId);
    }

    function _payout(address asset, address to, uint256 amount) internal {
        if (amount == 0) return;
        if (asset == address(0)) {
            (bool ok,) = to.call{value: amount}("");
            if (!ok) revert TransferFailed();
        } else {
            IERC20(asset).safeTransfer(to, amount);
        }
    }
}
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ProxyAdmin} from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import {TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import {ITransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {AuctionNFT} from "../src/AuctionNFT.sol";
import {AuctionHouseUUPS} from "../src/AuctionHouseUUPS.sol";
import {AuctionHouseCore} from "../src/AuctionHouseCore.sol";
import {AuctionHouseTransparent} from "../src/AuctionHouseTransparent.sol";
import {AuctionHouseV2UUPS} from "../src/AuctionHouseV2UUPS.sol";
import {AuctionHouseV2Transparent} from "../src/AuctionHouseV2Transparent.sol";
import {MockV3Aggregator} from "./mocks/MockV3Aggregator.sol";
import {MockERC20} from "./mocks/MockERC20.sol";

contract NftAuctionTest is Test {
    AuctionNFT nft;
    AuctionHouseUUPS house;
    MockV3Aggregator ethFeed;
    MockV3Aggregator tokenFeed;
    MockERC20 bidToken;

    address owner = makeAddr("owner");
    address seller = makeAddr("seller");
    address bidderEth = makeAddr("bidderEth");
    address bidderToken = makeAddr("bidderToken");

    function setUp() public {
        ethFeed = new MockV3Aggregator(8, 2000e8);
        tokenFeed = new MockV3Aggregator(8, 1e8);
        bidToken = new MockERC20();

        nft = new AuctionNFT(owner);
        AuctionHouseUUPS implementation = new AuctionHouseUUPS();
        ERC1967Proxy proxy = new ERC1967Proxy(
            address(implementation), abi.encodeCall(AuctionHouseUUPS.initialize, (owner, address(ethFeed), owner))
        );
        house = AuctionHouseUUPS(payable(address(proxy)));

        vm.prank(owner);
        house.setPriceFeed(address(bidToken), address(tokenFeed));

        vm.deal(bidderEth, 10 ether);
        vm.deal(bidderToken, 1 ether);
        bidToken.mint(bidderToken, 10_000 ether);
    }

    function test_mintAndTransfer() public {
        vm.prank(owner);
        uint256 tokenId = nft.mint(seller);
        assertEq(nft.ownerOf(tokenId), seller);

        vm.prank(seller);
        nft.transferFrom(seller, bidderEth, tokenId);
        assertEq(nft.ownerOf(tokenId), bidderEth);
    }

    function test_quoteUsd_ethAndErc20() public view {
        assertEq(house.quoteUsd(address(0), 1 ether), 2000 ether);
        assertEq(house.quoteUsd(address(bidToken), 2500 ether), 2500 ether);
    }

    function test_feeBpsTiers() public view {
        assertEq(house.feeBps(100 ether - 1), 250);
        assertEq(house.feeBps(100 ether), 150);
        assertEq(house.feeBps(1000 ether), 50);
    }

    function test_ethBidAndSettle_chargesDynamicFee() public {
        uint256 tokenId = _mintToSeller();
        uint256 auctionId = _list(tokenId, 1 days);

        vm.prank(bidderEth);
        house.bidWithEth{value: 0.01 ether}(auctionId);

        vm.warp(block.timestamp + 1 days);
        vm.prank(seller);
        house.settle(auctionId);

        uint256 fee = (0.01 ether * 250) / 10_000;
        assertEq(nft.ownerOf(tokenId), bidderEth);
        assertEq(owner.balance, fee);
        assertEq(seller.balance, 0.01 ether - fee);
        assertEq(address(house).balance, 0);
    }

    function test_erc20OutbidsEth_thenSettles() public {
        uint256 tokenId = _mintToSeller();
        uint256 auctionId = _list(tokenId, 1 days);

        vm.prank(bidderEth);
        house.bidWithEth{value: 1 ether}(auctionId);
        assertEq(bidderEth.balance, 9 ether);

        vm.startPrank(bidderToken);
        bidToken.approve(address(house), 2500 ether);
        house.bidWithErc20(auctionId, address(bidToken), 2500 ether);
        vm.stopPrank();

        assertEq(bidderEth.balance, 10 ether);
        assertEq(bidToken.balanceOf(address(house)), 2500 ether);

        uint256 lowUsd = house.quoteUsd(address(0), 1 ether);
        vm.prank(bidderEth);
        vm.expectRevert(abi.encodeWithSelector(AuctionHouseCore.BidTooLow.selector, lowUsd, 2500 ether));
        house.bidWithEth{value: 1 ether}(auctionId);

        vm.warp(block.timestamp + 1 days);
        house.settle(auctionId);

        uint256 fee = (2500 ether * 50) / 10_000;
        assertEq(nft.ownerOf(tokenId), bidderToken);
        assertEq(bidToken.balanceOf(owner), fee);
        assertEq(bidToken.balanceOf(seller), 2500 ether - fee);
        assertEq(bidToken.balanceOf(address(house)), 0);
    }

    function test_noBids_returnsNftToSeller() public {
        uint256 tokenId = _mintToSeller();
        uint256 auctionId = _list(tokenId, 1 hours);
        assertEq(nft.ownerOf(tokenId), address(house));

        vm.warp(block.timestamp + 1 hours);
        house.settle(auctionId);
        assertEq(nft.ownerOf(tokenId), seller);
    }

    function test_settleBeforeEndReverts() public {
        uint256 tokenId = _mintToSeller();
        uint256 auctionId = _list(tokenId, 1 days);
        vm.expectRevert(abi.encodeWithSelector(AuctionHouseCore.AuctionOngoing.selector, auctionId));
        house.settle(auctionId);
    }

    function test_missingPriceFeedReverts() public {
        MockERC20 other = new MockERC20();
        vm.expectRevert(abi.encodeWithSelector(AuctionHouseCore.MissingPriceFeed.selector, address(other)));
        house.quoteUsd(address(other), 1 ether);
    }

    function test_transparentUpgrade_preservesAuctionAndAllowsV2() public {
        AuctionHouseTransparent implementation = new AuctionHouseTransparent();
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(implementation),
            owner,
            abi.encodeCall(AuctionHouseTransparent.initialize, (owner, address(ethFeed), owner))
        );
        bytes32 adminSlot = 0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103;
        ProxyAdmin admin = ProxyAdmin(address(uint160(uint256(vm.load(address(proxy), adminSlot)))));
        AuctionHouseTransparent houseT = AuctionHouseTransparent(payable(address(proxy)));

        vm.prank(owner);
        houseT.setPriceFeed(address(bidToken), address(tokenFeed));

        uint256 tokenId = _mintToSeller();
        uint256 auctionId = _listOn(address(houseT), tokenId, 2 days);
        vm.prank(bidderEth);
        houseT.bidWithEth{value: 1 ether}(auctionId);
        assertEq(houseT.implementationVersion(), 1);

        AuctionHouseV2Transparent next = new AuctionHouseV2Transparent();
        vm.prank(bidderEth);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, bidderEth));
        admin.upgradeAndCall(ITransparentUpgradeableProxy(address(proxy)), address(next), "");

        vm.prank(owner);
        admin.upgradeAndCall(ITransparentUpgradeableProxy(address(proxy)), address(next), "");

        AuctionHouseV2Transparent upgraded = AuctionHouseV2Transparent(payable(address(proxy)));
        assertEq(upgraded.implementationVersion(), 2);
        (address listedSeller,, uint256 listedTokenId,,,,,,) = upgraded.auctions(auctionId);
        assertEq(listedSeller, seller);
        assertEq(listedTokenId, tokenId);

        upgraded.pingAfterUpgrade();
        assertEq(upgraded.postUpgradePings(), 1);
    }

    function test_upgrade_preservesAuctionAndAllowsV2() public {
        uint256 tokenId = _mintToSeller();
        uint256 auctionId = _list(tokenId, 2 days);
        vm.prank(bidderEth);
        house.bidWithEth{value: 1 ether}(auctionId);
        assertEq(house.implementationVersion(), 1);

        AuctionHouseV2UUPS next = new AuctionHouseV2UUPS();
        vm.prank(bidderEth);
        vm.expectRevert(abi.encodeWithSelector(OwnableUpgradeable.OwnableUnauthorizedAccount.selector, bidderEth));
        house.upgradeToAndCall(address(next), "");

        vm.prank(owner);
        house.upgradeToAndCall(address(next), "");

        AuctionHouseV2UUPS upgraded = AuctionHouseV2UUPS(payable(address(house)));
        assertEq(upgraded.implementationVersion(), 2);
        (address listedSeller,, uint256 listedTokenId,,,,,,) = upgraded.auctions(auctionId);
        assertEq(listedSeller, seller);
        assertEq(listedTokenId, tokenId);
        assertEq(upgraded.quoteUsd(address(0), 1 ether), 2000 ether);

        upgraded.pingAfterUpgrade();
        assertEq(upgraded.postUpgradePings(), 1);

        vm.warp(block.timestamp + 2 days);
        upgraded.settle(auctionId);
        uint256 fee = (1 ether * 50) / 10_000;
        assertEq(nft.ownerOf(tokenId), bidderEth);
        assertEq(seller.balance, 1 ether - fee);
    }

    function test_adminGuardsAndBadPrice() public {
        vm.prank(seller);
        vm.expectRevert(abi.encodeWithSelector(OwnableUpgradeable.OwnableUnauthorizedAccount.selector, seller));
        house.setPriceFeed(address(bidToken), address(ethFeed));

        vm.prank(owner);
        house.setFeeRecipient(seller);
        assertEq(house.feeRecipient(), seller);

        vm.prank(owner);
        vm.expectRevert(AuctionHouseCore.ZeroAddress.selector);
        house.setFeeRecipient(address(0));

        ethFeed.setAnswer(0);
        vm.expectRevert(abi.encodeWithSelector(AuctionHouseCore.BadPrice.selector, address(ethFeed)));
        house.quoteUsd(address(0), 1 ether);
    }

    function test_zeroDurationReverts() public {
        uint256 tokenId = _mintToSeller();
        vm.startPrank(seller);
        nft.approve(address(house), tokenId);
        vm.expectRevert(AuctionHouseCore.InvalidDuration.selector);
        house.createAuction(address(nft), tokenId, 0);
        vm.stopPrank();
    }

    function test_cannotReinitialize() public {
        vm.expectRevert(Initializable.InvalidInitialization.selector);
        house.initialize(owner, address(ethFeed), owner);

        AuctionHouseUUPS implementation = new AuctionHouseUUPS();
        vm.expectRevert(Initializable.InvalidInitialization.selector);
        implementation.initialize(owner, address(ethFeed), owner);
    }

    function _mintToSeller() internal returns (uint256 tokenId) {
        vm.prank(owner);
        tokenId = nft.mint(seller);
    }

    function _list(uint256 tokenId, uint256 duration) internal returns (uint256 auctionId) {
        return _listOn(address(house), tokenId, duration);
    }

    function _listOn(address houseAddr, uint256 tokenId, uint256 duration) internal returns (uint256 auctionId) {
        vm.startPrank(seller);
        nft.approve(houseAddr, tokenId);
        auctionId = AuctionHouseCore(houseAddr).createAuction(address(nft), tokenId, duration);
        vm.stopPrank();
    }
}

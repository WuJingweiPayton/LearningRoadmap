// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {AuctionHouseCore} from "./AuctionHouseCore.sol";

/// @title NFT auction house (transparent proxy). Deploy behind `TransparentUpgradeableProxy`.
contract AuctionHouseTransparent is AuctionHouseCore {
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(address initialOwner, address ethUsdFeed, address feeRecipient_) public initializer {
        _initializeAuctionHouse(initialOwner, ethUsdFeed, feeRecipient_);
    }
}
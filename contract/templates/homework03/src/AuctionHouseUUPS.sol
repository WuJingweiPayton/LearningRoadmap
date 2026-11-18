// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {AuctionHouseCore} from "./AuctionHouseCore.sol";

/// @title NFT auction house (UUPS). Deploy behind `ERC1967Proxy` and call `initialize`.
contract AuctionHouseUUPS is AuctionHouseCore, UUPSUpgradeable {
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(address initialOwner, address ethUsdFeed, address feeRecipient_) public initializer {
        _initializeAuctionHouse(initialOwner, ethUsdFeed, feeRecipient_);
        __UUPSUpgradeable_init();
    }

    function _authorizeUpgrade(address) internal override onlyOwner {}
}
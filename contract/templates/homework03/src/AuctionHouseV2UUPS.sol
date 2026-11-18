// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {AuctionHouseUUPS} from "./AuctionHouseUUPS.sol";

/// @title UUPS auction house upgrade. Storage appended only; existing auctions stay valid.
contract AuctionHouseV2UUPS is AuctionHouseUUPS {
    uint256 public postUpgradePings;

    function implementationVersion() public pure override returns (uint256) {
        return 2;
    }

    function pingAfterUpgrade() external {
        postUpgradePings += 1;
    }
}
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {AuctionHouseTransparent} from "./AuctionHouseTransparent.sol";

/// @title Transparent-proxy auction house upgrade. Storage appended only.
contract AuctionHouseV2Transparent is AuctionHouseTransparent {
    uint256 public postUpgradePings;

    function implementationVersion() public pure override returns (uint256) {
        return 2;
    }

    function pingAfterUpgrade() external {
        postUpgradePings += 1;
    }
}
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Script, console2} from "forge-std/Script.sol";
import {AuctionHouseUUPS} from "../src/AuctionHouseUUPS.sol";
import {AuctionHouseV2UUPS} from "../src/AuctionHouseV2UUPS.sol";

/// @notice UUPS: upgrade via `upgradeToAndCall` on the proxy. Caller must be the proxy owner.
///         Set `AUCTION_PROXY` from `DeployAuctionUUPSScript`.
contract UpgradeAuctionUUPSScript is Script {
    function run() external {
        address proxy = vm.envAddress("AUCTION_PROXY");

        vm.startBroadcast();
        AuctionHouseV2UUPS implementation = new AuctionHouseV2UUPS();
        AuctionHouseUUPS(payable(proxy)).upgradeToAndCall(address(implementation), "");
        vm.stopBroadcast();

        console2.log("AuctionHouseV2UUPS implementation", address(implementation));
        console2.log("proxy", proxy);
    }
}
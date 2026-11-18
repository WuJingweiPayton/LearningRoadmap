// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Script, console2} from "forge-std/Script.sol";
import {ProxyAdmin} from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import {ITransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import {AuctionHouseV2Transparent} from "../src/AuctionHouseV2Transparent.sol";

/// @notice Transparent proxy: upgrade through `ProxyAdmin`. Caller must own the admin.
///         Set `AUCTION_PROXY` and `PROXY_ADMIN` from `DeployAuctionTransparentScript`.
contract UpgradeAuctionTransparentScript is Script {
    function run() external {
        address proxy = vm.envAddress("AUCTION_PROXY");
        address proxyAdmin = vm.envAddress("PROXY_ADMIN");

        vm.startBroadcast();
        AuctionHouseV2Transparent implementation = new AuctionHouseV2Transparent();
        ProxyAdmin(proxyAdmin).upgradeAndCall(ITransparentUpgradeableProxy(proxy), address(implementation), "");
        vm.stopBroadcast();

        console2.log("AuctionHouseV2Transparent implementation", address(implementation));
        console2.log("proxy", proxy);
        console2.log("ProxyAdmin", proxyAdmin);
    }
}
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Script, console2} from "forge-std/Script.sol";
import {TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import {AuctionNFT} from "../src/AuctionNFT.sol";
import {AuctionHouseTransparent} from "../src/AuctionHouseTransparent.sol";

/// @notice Deploy NFT + transparent proxy auction. Requires `ETH_USD_FEED`.
contract DeployAuctionTransparentScript is Script {
    function run() external {
        address ethUsdFeed = vm.envAddress("ETH_USD_FEED");

        vm.startBroadcast();
        address deployer = msg.sender;
        AuctionNFT nft = new AuctionNFT(deployer);
        AuctionHouseTransparent implementation = new AuctionHouseTransparent();
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(implementation),
            deployer,
            abi.encodeCall(AuctionHouseTransparent.initialize, (deployer, ethUsdFeed, deployer))
        );
        vm.stopBroadcast();

        bytes32 adminSlot = 0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103;
        address proxyAdmin = address(uint160(uint256(vm.load(address(proxy), adminSlot))));

        console2.log("AuctionNFT", address(nft));
        console2.log("AuctionHouseTransparent implementation", address(implementation));
        console2.log("ProxyAdmin", proxyAdmin);
        console2.log("AuctionHouse transparent proxy", address(proxy));
    }
}
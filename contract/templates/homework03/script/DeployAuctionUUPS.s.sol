// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Script, console2} from "forge-std/Script.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {AuctionNFT} from "../src/AuctionNFT.sol";
import {AuctionHouseUUPS} from "../src/AuctionHouseUUPS.sol";

/// @notice UUPS: deploy NFT + `ERC1967Proxy`. Requires `ETH_USD_FEED` (Chainlink aggregator).
contract DeployAuctionUUPSScript is Script {
    function run() external {
        address ethUsdFeed = vm.envAddress("ETH_USD_FEED");

        vm.startBroadcast();
        address deployer = msg.sender;
        AuctionNFT nft = new AuctionNFT(deployer);
        AuctionHouseUUPS implementation = new AuctionHouseUUPS();
        ERC1967Proxy proxy = new ERC1967Proxy(
            address(implementation), abi.encodeCall(AuctionHouseUUPS.initialize, (deployer, ethUsdFeed, deployer))
        );
        vm.stopBroadcast();

        console2.log("AuctionNFT", address(nft));
        console2.log("AuctionHouseUUPS implementation", address(implementation));
        console2.log("AuctionHouseUUPS proxy", address(proxy));
    }
}
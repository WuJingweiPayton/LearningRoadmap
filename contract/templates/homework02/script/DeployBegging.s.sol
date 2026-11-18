// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Script} from "forge-std/Script.sol";
import {BeggingContract} from "../src/BeggingContract.sol";

contract DeployBeggingScript is Script {
    function run() external returns (BeggingContract) {
        vm.startBroadcast();
        BeggingContract begging = new BeggingContract();
        vm.stopBroadcast();
        return begging;
    }
}

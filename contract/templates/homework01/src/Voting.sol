// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

/// @title Simple voting contract (homework 01)
contract Voting {
    mapping(string => uint256) public votes;
    string[] private _candidates;

    function vote(string calldata candidate) external {
        if (!_candidateTracked(candidate)) {
            _candidates.push(candidate);
        }
        votes[candidate]++;
    }

    function getVotes(string calldata candidate) external view returns (uint256) {
        return votes[candidate];
    }

    function resetVotes() external {
        for (uint256 i = 0; i < _candidates.length; i++) {
            delete votes[_candidates[i]];
        }
    }

    function _candidateTracked(string calldata candidate) private view returns (bool) {
        bytes32 key = keccak256(bytes(candidate));
        for (uint256 i = 0; i < _candidates.length; i++) {
            if (keccak256(bytes(_candidates[i])) == key) {
                return true;
            }
        }
        return false;
    }
}

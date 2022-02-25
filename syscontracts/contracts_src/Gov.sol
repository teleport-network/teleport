// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

// generate go: abigen --sol Gov.sol --pkg gov --out generated.go

contract Gov {
    event Voted(address voter, uint64 proposalID, uint32 voteOption);
    event VotedWeighted(
        address voter,
        uint64 proposalID,
        OptionWeight[] options
    );

    struct OptionWeight {
        uint32 option;
        uint64 weight; // 2 decimal place, e.g. 20 = 20%, 80 = 80%
    }

    function vote(uint64 proposalID, uint32 voteOption) external {
        emit Voted(msg.sender, proposalID, voteOption);
    }

    function vote(uint64 proposalID, OptionWeight[] memory options) external {
        emit VotedWeighted(msg.sender, proposalID, options);
    }
}

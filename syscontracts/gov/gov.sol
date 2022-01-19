// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.6.8;
pragma experimental ABIEncoderV2;

contract Gov {
    event Voted(address voter, uint64 proposalId, uint32 voteOption);
    event VotedWeighted(
        address voter,
        uint64 proposalId,
        OptionWeight[] options
    );

    struct OptionWeight {
        uint32 option;
        uint64 weight; // 2 decimal place, e.g. 20 = 20%, 80 = 80%
    }

    function vote(uint64 proposalId, uint32 voteOption) external {
        emit Voted(msg.sender, proposalId, voteOption);
    }

    function vote(uint64 proposalId, OptionWeight[] calldata options) external {
        emit VotedWeighted(msg.sender, proposalId, options);
    }
}

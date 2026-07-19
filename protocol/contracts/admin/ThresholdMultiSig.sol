// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity 0.8.24;

contract ThresholdMultiSig {
    // Defines who can approve and how many approvals are required
    address[] private owners;
    uint256 public threshold;
    mapping(address => bool) public isOwner;

    // Tracks the approval status of each transaction.
    mapping(bytes32 => uint256) public approvalCount; // tx hash => number of approvals
    mapping(bytes32 => bool) public executed; // tx hash => executed
    mapping(bytes32 => mapping(address => bool)) public hasApproved; // tx hash => owner => approved

    // Events
    event TransactionApproved(bytes32 indexed txHash, address indexed owner, uint256 approvalCount);
    event TransactionExecuted(bytes32 indexed txHash, address indexed executor, address indexed target);

    // Modifiers
    modifier onlyOwner() {
        require(isOwner[msg.sender], "Not an owner");
        _;
    }

    constructor(address[] memory owners_, uint256 threshold_) {
        require(owners_.length > 0, "Owners required");
        require(threshold_ > 0 && threshold_ <= owners_.length, "Invalid threshold");

        for (uint256 i = 0; i < owners_.length; i++) {
            address owner = owners_[i];
            require(owner != address(0), "Invalid owner");
            require(!isOwner[owner], "Owner not unique");

            isOwner[owner] = true;
            owners.push(owner);
        }

        threshold = threshold_;
    }

    function ownerCount() external view returns (uint256) {
        return owners.length;
    }

    function getOwner(uint256 index) external view returns (address) {
        require(index < owners.length, "Index out of bounds");
        return owners[index];
    }

    function getTransactionHash(address target, uint256 value, bytes calldata data, uint256 nonce)
        public
        view
        returns (bytes32)
    {
        return keccak256(abi.encode(address(this), block.chainid, target, value, keccak256(data), nonce));
    }

    function approveTransaction(address target, uint256 value, bytes calldata data, uint256 nonce)
        external
        onlyOwner
        returns (bytes32 txHash) 
    {
        require(target != address(0), "Invalid target");

        txHash = getTransactionHash(target, value, data, nonce);
        require(!executed[txHash], "Transaction already executed");
        require(!hasApproved[txHash][msg.sender], "Already approved by this owner");

        hasApproved[txHash][msg.sender] = true;
        approvalCount[txHash] += 1;

        emit TransactionApproved(txHash, msg.sender, approvalCount[txHash]);
    }

    function executeTransaction(address target, uint256 value, bytes calldata data, uint256 nonce)
        external
        onlyOwner
        returns (bytes memory result)
    {
        bytes32 txHash = getTransactionHash(target, value, data, nonce);
        require(!executed[txHash], "Transaction already executed");
        require(approvalCount[txHash] >= threshold, "Not enough approvals");
        
        // Execute the transaction follow CEI (Check, Effect, Interaction) pattern 
        executed[txHash] = true;
        bool success;
        (success, result) = target.call{value: value}(data);
        require(success, "Transaction execution failed");

        emit TransactionExecuted(txHash, msg.sender, target);
    }

    // receive() runs when ETH arrives with empty calldata
    receive() external payable {}
}
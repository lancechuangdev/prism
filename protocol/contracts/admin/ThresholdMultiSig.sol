// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;

contract ThresholdMultiSig {
    // Defines who can approve and how many approvals are required
    address[] private owners;
    uint256 public threshold;
    uint256 public configurationVersion;
    mapping(address => bool) public isOwner;

    // Tracks the approval status of each transaction.
    mapping(bytes32 => uint256) public approvalCount; // tx hash => number of approvals
    mapping(bytes32 => uint256) public transactionConfigurationVersion;
    mapping(bytes32 => bool) public executed; // tx hash => executed
    mapping(bytes32 => mapping(address => bool)) public hasApproved; // tx hash => owner => approved

    // Events
    event TransactionApproved(bytes32 indexed txHash, address indexed owner, uint256 approvalCount);
    event TransactionExecuted(bytes32 indexed txHash, address indexed executor, address indexed target);
    event OwnerAdded(address indexed owner);
    event OwnerRemoved(address indexed owner);
    event OwnerReplaced(address indexed oldOwner, address indexed newOwner);
    event ThresholdChanged(uint256 oldThreshold, uint256 newThreshold);

    // Modifiers
    modifier onlyOwner() {
        require(isOwner[msg.sender], "Not an owner");
        _;
    }

    modifier onlySelf() {
        require(msg.sender == address(this), "Only multisig");
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

    function addOwner(address owner) external onlySelf {
        require(owner != address(0), "Invalid owner");
        require(!isOwner[owner], "Owner not unique");

        isOwner[owner] = true;
        owners.push(owner);
        configurationVersion += 1;

        emit OwnerAdded(owner);
    }

    function removeOwner(address owner) external onlySelf {
        require(isOwner[owner], "Not an owner");
        require(owners.length > 1, "Owners required");
        require(threshold <= owners.length - 1, "Threshold exceeds owner count");

        uint256 ownerIndex = _ownerIndex(owner);
        uint256 lastIndex = owners.length - 1;
        if (ownerIndex != lastIndex) {
            owners[ownerIndex] = owners[lastIndex];
        }
        owners.pop();
        isOwner[owner] = false;
        configurationVersion += 1;

        emit OwnerRemoved(owner);
    }

    function replaceOwner(address oldOwner, address newOwner) external onlySelf {
        require(isOwner[oldOwner], "Not an owner");
        require(newOwner != address(0), "Invalid owner");
        require(!isOwner[newOwner], "Owner not unique");

        owners[_ownerIndex(oldOwner)] = newOwner;
        isOwner[oldOwner] = false;
        isOwner[newOwner] = true;
        configurationVersion += 1;

        emit OwnerReplaced(oldOwner, newOwner);
    }

    function changeThreshold(uint256 newThreshold) external onlySelf {
        require(newThreshold > 0 && newThreshold <= owners.length, "Invalid threshold");

        uint256 oldThreshold = threshold;
        threshold = newThreshold;
        configurationVersion += 1;

        emit ThresholdChanged(oldThreshold, newThreshold);
    }

    function getTransactionHash(address target, uint256 value, bytes calldata data, uint256 nonce)
        public
        view
        returns (bytes32)
    {
        return keccak256(
            abi.encode(address(this), block.chainid, configurationVersion, target, value, keccak256(data), nonce)
        );
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

        if (approvalCount[txHash] == 0) {
            // Store version + 1 so zero means this hash has never received an approval.
            transactionConfigurationVersion[txHash] = configurationVersion + 1;
        }
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

    function _ownerIndex(address owner) private view returns (uint256) {
        for (uint256 i = 0; i < owners.length; i++) {
            if (owners[i] == owner) {
                return i;
            }
        }
        revert("Not an owner");
    }

    // receive() runs when ETH arrives with empty calldata
    receive() external payable {}
}

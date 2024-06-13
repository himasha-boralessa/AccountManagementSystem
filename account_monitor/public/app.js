document.addEventListener('DOMContentLoaded', (event) => {
    fetchTransactions();
  });
  
  async function fetchTransactions() {
    try {
      const response = await fetch('http://localhost:3001/transactions');
      const data = await response.json();
      displayTransactions(data.transactions);
      displayBalance(data.balance);
    } catch (error) {
      console.error('Error fetching transactions:', error);
    }
  }
  
  function displayTransactions(transactions) {
    const transactionList = document.getElementById('transaction-list');
    transactionList.innerHTML = ''; // Clear any existing content
    transactions.forEach(transaction => {
      const listItem = document.createElement('li');
      listItem.textContent = `${transaction.time}: ${transaction.amount}`;
      transactionList.appendChild(listItem);
    });
  }
  
  function displayBalance(balance) {
    const balanceElement = document.getElementById('balance');
    balanceElement.textContent = `Balance: ${balance}`;
  }
  
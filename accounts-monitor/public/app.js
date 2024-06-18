document.addEventListener('DOMContentLoaded', (event) => {
    fetchTransactions();
  });
  
  async function fetchTransactions() {
    try {
      const response = await fetch('/monitor');
      const data = await response.json();
      displayTransactions(data.transactions);
      // displayBalance(data.balance);
    } catch (error) {
      console.error('Error fetching transactions:', error);
    }
  }
  
  function displayTransactions(transactions) {
    const transactionList = document.getElementById('transaction-list');
    transactionList.innerHTML = ''; // Clear any existing content
    transactions.forEach(transaction => {
        const row = document.createElement('tr');
        const clientCell = document.createElement('td');
        const amountCell = document.createElement('td');
        const balanceCell = document.createElement('td');
        
        clientCell.textContent = transaction.client_id;
        amountCell.textContent = transaction.amount;
        balanceCell.textContent = transaction.balance;
        
        row.appendChild(clientCell);
        row.appendChild(amountCell);
        row.appendChild(balanceCell);

        transactionList.appendChild(row);
    });
  }
  
  // function displayBalance(balance) {
  //   const balanceElement = document.getElementById('balance');
  //   balanceElement.textContent = `Balance: ${balance}`;
  // }
  
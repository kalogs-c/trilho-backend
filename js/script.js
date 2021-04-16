const transactionsUL = document.querySelector('#transactions')
const totalElement = document.querySelector('#balance')
const incomeElement = document.querySelector('#money-plus')
const expenseElement = document.querySelector('#money-minus')
const form = document.querySelector('#form')
const inputTransactionName = document.querySelector('#text')
const inputTransactionAmount = document.querySelector('#amount')
const expButton = document.querySelector('#exp-button')
const incButton = document.querySelector('#inc-button')


const localStorageTransactions = JSON.parse(localStorage.getItem('transactions'))
let transactions = localStorage.getItem('transactions') !== null ? localStorageTransactions : []

const removeTransaction = ID => {
    transactions = transactions.filter(transaction => transaction.id !== ID)
    init()
}

const addTransactionIntoDOM = transaction => {
    const operator = transaction.amount < 0 ? '-' : '+'
    const CSSclass = transaction.amount < 0 ? 'minus' : 'plus'
    const absoluteValue_amount = Math.abs(transaction.amount)
    const li = document.createElement('li')

    li.classList.add(CSSclass)
    li.innerHTML = `
        ${transaction.name} <span>${operator} $ ${absoluteValue_amount}</span>
        <button class="delete-btn" onclick="removeTransaction(${transaction.id})">x</button>
    `

    transactionsUL.prepend(li)
}

const updateBalanceValues = () => {
    const transactionAmounts = transactions
        .map(transaction => transaction.amount)

    const total = transactionAmounts
        .reduce((accumulator, transaction) => accumulator + transaction, 0)
        .toFixed(2)

    const income = transactionAmounts
        .filter(value => value > 0)
        .reduce((accumulator, value) => accumulator + value, 0)
        .toFixed(2)

    const expense = Math.abs(transactionAmounts
        .filter(value => value < 0)
        .reduce((accumulator, value) => accumulator + value, 0))
        .toFixed(2)

    expenseElement.textContent = `$ ${expense}`
    incomeElement.textContent = `$ ${income}`
    totalElement.textContent = `$ ${total}`
}

const init = () => {
    transactionsUL.innerHTML = ''
    transactions.forEach(addTransactionIntoDOM)
    updateBalanceValues()
}

init()

const updateLocalStorage = () => {
    localStorage.setItem('transactions', JSON.stringify(transactions))
}

const generateID = () => Math.round(Math.random() * 1000)

form.addEventListener('submit', event => {

})

expButton.addEventListener('click', () => {
    const transactionName = inputTransactionName.value.trim()
    let transactionAmount = inputTransactionAmount.value.trim()

    if (transactionName === '' || transactionAmount === '') {
        alert('Por favor, preencha tanto o nome quanto o valor da transação')
        return
    }

    const transaction = {
        id: generateID(),
        name: transactionName,
        amount: Number(transactionAmount * (-1))
    }

    transactions.push(transaction)
    init()
    updateLocalStorage()

    inputTransactionName.value = ''
    inputTransactionAmount.value = ''
})

incButton.addEventListener('click', () => {
    const transactionName = inputTransactionName.value.trim()
    let transactionAmount = inputTransactionAmount.value.trim()

    if (transactionName === '' || transactionAmount === '') {
        alert('Por favor, preencha tanto o nome quanto o valor da transação')
        return
    }

    const transaction = {
        id: generateID(),
        name: transactionName,
        amount: Number(transactionAmount)
    }

    transactions.push(transaction)
    init()
    updateLocalStorage()

    inputTransactionName.value = ''
    inputTransactionAmount.value = ''
})
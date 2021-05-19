const formDiv = document.querySelector('.form')
const container = document.querySelector('.container')
const button = document.querySelector('.open-add_button')
const closeButton = document.querySelector('.close-button')


const showAddPage = () => {
    container.classList.add('blur')
    formDiv.style.display = 'block'
}

const hideAddPage = () => {
    container.classList.remove('blur')
    formDiv.style.display = 'none'
}

button.addEventListener('click', showAddPage)
closeButton.addEventListener('click', hideAddPage)

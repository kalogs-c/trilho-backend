const formDiv = document.querySelector('.form')
const container = document.querySelector('.container')
const button = document.querySelector('.open-add_button')

button.addEventListener('click',() => {
    container.classList.add('blur')
    formDiv.style.display = 'block'
})

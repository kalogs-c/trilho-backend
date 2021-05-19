const formDiv = document.querySelector('.form')
const container = document.querySelector('.container')
const button = document.querySelector('.open-add_button')
const closeButton = document.querySelector('.close-button')

const showAddPage = () => {
    container.classList.add('blur')
    formDiv.style.display = 'block'

    const plusIcon = document.querySelector('.open-add_button > .fa-plus')
    plusIcon.classList.add('rotate')
    
    if (formDiv.classList.contains('show')) {
        hideAddPage()
        return
    }

    formDiv.classList.add('show')
}

const hideAddPage = () => {
    container.classList.remove('blur')
    formDiv.style.display = 'none'
    formDiv.classList.remove('show')
    const plusIcon = document.querySelector('.open-add_button > .fa-plus')
    plusIcon.classList.remove('rotate')
}

button.addEventListener('click', showAddPage)
closeButton.addEventListener('click', hideAddPage)

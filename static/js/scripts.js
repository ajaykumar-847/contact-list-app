
// Function to generate contact editing form with the existing data
function editContact(name, email, phone_number, address, company_name) {
    // Set values for each field in the edit form
    document.getElementById('edit_name').value = name;
    document.getElementById('edit_name_new').value = name;
    document.getElementById('edit_email').value = email;
    document.getElementById('edit_phone_number').value = phone_number;
    document.getElementById('edit_address').value = address;
    document.getElementById('edit_company_name').value = company_name;

    // Display the edit form
    document.getElementById('edit_form').style.display = 'block';
}

// Function to delete a contact
function deleteContact(name, username) {
    // Confirm deletion with the user
    if (confirm('Are you sure you want to delete this contact?')) {
        // console.log("User = ",username);
        // Create a form dynamically
        const form = document.createElement('form');
        form.method = 'post';
        form.action = '/deleteContact';
        
        // Create hidden input field for username
        const usernameField = document.createElement('input');
        usernameField.type = 'hidden';
        usernameField.name = 'username';
        usernameField.value = username; 

        // Create hidden input field for contact name
        const nameField = document.createElement('input');
        nameField.type = 'hidden';
        nameField.name = 'name';
        nameField.value = name;

        // Append the input fields to the form
        form.appendChild(usernameField);
        form.appendChild(nameField);

        // Append form to the document body and submit it
        document.body.appendChild(form);
        form.submit();
    }
}

// Event listener for when the DOM content is loaded
document.addEventListener('DOMContentLoaded', function() {
    // Function to validate signup form submission
    function validateSignupForm(event) {
        // Get values from the signup form fields
        const username = document.getElementById('username').value;
        const phoneNumber = document.getElementById('phone_number').value;
        const password = document.getElementById('password').value;

        /* Check if any field is empty, if there are empty fields 
        prevent form submission and display alert message*/
        if (username === '' || phoneNumber === '' || password === '') {
            alert('All fields are required.');
            event.preventDefault();
        }
    }

    // Function to validate login form submission
    function validateLoginForm(event) {
        // Get values from the login form fields
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;

        /* Check if any field is empty, if there are empty fields 
        prevent form submission and display alert message*/
        if (username === '' || password === '') {
            alert('Both fields are required.');
            event.preventDefault();
        }
    }

    // Event listener for form submission on signup form
    const signupForm = document.querySelector('form[action="/signup"]');
    if (signupForm) {
        signupForm.addEventListener('submit', validateSignupForm);
    }

    // Add event listener for form submission on login form
    const loginForm = document.querySelector('form[action="/login"]');
    if (loginForm) {
        loginForm.addEventListener('submit', validateLoginForm);
    }

    // Function to validate contacts form submission
    function validateContactsForm(event) {
        // Get values from the contacts form fields
        const name = document.getElementById('name').value;
        const email = document.getElementById('email').value;
        const phoneNumber = document.getElementById('phone_number').value;
        const address = document.getElementById('address').value;
        const companyName = document.getElementById('company_name').value;

        /* Check if any field is empty, if there are empty fields 
        prevent form submission and display alert message*/
        if (name === '' || email === '' || phoneNumber === '' || address === '' || companyName === '') {
            alert('All fields are required.');
            event.preventDefault();
        }

    }

    // Add event listener for form submission on contacts form
    const contactsForm = document.querySelector('form[action="/contacts"]');
    if (contactsForm) {
        contactsForm.addEventListener('submit', validateContactsForm);
    }
});






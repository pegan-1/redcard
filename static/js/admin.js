/*
admin.js
Scripting for the admin panel.

@author  Peter Egan
@lastUpdated 2021-10-10

Copyright (c) 2021 kiercam llc
*/

// Document Elements
const editor = document.getElementById('editor');
const blogSummary = document.getElementById('blog_summary');
const blogTitle = document.getElementById('blog_title');

const publishButton = document.getElementById('publish');

//add the toolbar options
let myToolbar= [
    // [{ 'header': [1, 2, 3, 4, 5, 6, false] }],
    [{ 'size': ['small', false, 'large', 'huge'] }],  // custom dropdown
    ['bold', 'italic', 'underline', 'strike'],  
    [{ 'color': [] }, { 'background': [] }], 
    [{ 'list': 'ordered'}, { 'list': 'bullet' }], 
    ['blockquote', 'code-block'], 
    [{ 'indent': '-1'}, { 'indent': '+1' }],          // outdent/indent
    [{ 'align': [] }],         
    // [{ 'font': [] }],
    // [{ 'align': [] }],
    // ['link'],        
    ['image'], 
    ['clean']         
];

let quill = new Quill('#editor', {
    theme: 'snow',
    modules: {
        toolbar: {
            container: myToolbar,
            // handlers: {
            //     // handlers object will be merged with default handlers object
            //     'link': function(value) {
            //         if (value) {
            //             let urlLink = prompt('Enter the URL');
            //             this.quill.format('link', urlLink);
            //         } else {
            //             this.quill.format('link', false);
            //         }
            //     } 
            // }
        },
        // imageResize: {
        //     displaySize: true
        // },
    },
});

// Set up event listeners
publishButton.addEventListener('click',async _ => {
    // Grab the editor content.
    let content = quill.container.firstChild.innerHTML

    // Send the blog post to the server
    try{
        // Construct the payload
        console.log("The blog summary is ")
        console.log(blogSummary.value)
        let payload = {"title": blogTitle.value, "summary": blogSummary.value, "content":content};
 
        // Send the  HTTP Request
        let request = new XMLHttpRequest();
        request.open("POST", '/admin')
        request.setRequestHeader("Accept", "application/json");
        request.setRequestHeader("Content-Type", "application/json");
        request.send(JSON.stringify(payload))

        // Process the HTTP Reponse
        request.onload = function(){
            if (request.readyState === request.DONE) {
                if (request.status === 200) {
                    if(request.responseText === 'OK'){
                        window.location.replace('/blog')
                    }
                }
            }
        }
        //TODO - Better error handling.
    } catch(err) {
      // TBD... How do I manage errors?
      console.error(`Error: ${err}`);
    }
});
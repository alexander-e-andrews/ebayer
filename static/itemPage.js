
function updateItem() {

    //event.preventDefault()

    let f = document.getElementById("form")

    const XHR = new XMLHttpRequest(),
        FD = new FormData(f);

    // Define what happens in case of error
    XHR.addEventListener(' error', function (event) {
        alert('Oops! Something went wrong.');
    });

    // Set up our request
    XHR.open('POST', '');

    // Send our FormData object; HTTP headers are set automatically
    XHR.send(FD);
    
}

function uploadImages(){

    im = document.getElementById("images")
    for(var i = 0; i < im.files.length; i ++){
        console.log(im.files[i])
        sendFile(im.files[i])
    }


    return false
}

function sendFile(file){
    let formData = new FormData();
    let xhr = new XMLHttpRequest();

    formData.set("image", file)
    formData.set("itemID", document.getElementById("itemID").value)
    for (var value of formData.values()) {
        console.log(value); 
     }

    xhr.open("POST", "/upload")

    xhr.send(formData)
}

function nextPage(num){
    num = num + 1
    window.location.href = "/item/"+num
}
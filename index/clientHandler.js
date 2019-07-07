


function sendMessage()
{
    var xmlHttpRequest = new XMLHttpRequest();
    
    xmlHttpRequest.open('POST', 'sendMessage', false);
    var message = document.getElementById("messageWriter").value
    xmlHttpRequest.send(message);
    
    if(xmlHttpRequest.status != 200){
        alert("Сообщение не отправлено")
    }
    else
    {
        var chatField = document.getElementById("chatMessages")
        chatField.value += xmlHttpRequest.responseText
    }
}
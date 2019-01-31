runtime = setInterval(function() {
    for (var ii=1; ii<4; ii++) {
        $.ajax({
            url: '/check/' + ii,
            type: 'get',
            dataType: 'text',
            success: function (data) {
                document.getElementById("code_" + ii).innerText = data;
            },
            error: function (request, status, error) {   // 오류가 발생했을 때 호출된다.
                console.log("code:" + request.status + "\n" + "message:" + request.responseText + "\n" + "error:" + error);

            }


        })
    }
}, 1500);
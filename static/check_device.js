runtime = setInterval(function() {
    for (var ii=1; ii<4; ii++) {
        $.ajax({
            url: '/check/' + ii,
            type: 'get',
            dataType: 'text',
            async: true,
            success: function (result) {
                var data = JSON.parse(result);
                document.getElementById("code_" + data.code).innerText = data.code + " " + data.result;

                if (data.result == false) {
                    document.getElementById("code_" + data.code).style.backgroundColor = "#ff000f";
                }
                else if (data.result == true) {
                    document.getElementById("code_" + data.code).style.backgroundColor = "#00ff00";
                }
                else {
                    document.getElementById("code_" + data.code).style.backgroundColor = "#a0a0a0";
                }
            },
            error: function (request, status, error) {   // 오류가 발생했을 때 호출된다.
                console.log("code:" + request.status + "\n" + "message:" + request.responseText + "\n" + "error:" + error);
                document.getElementById("code_" + ii).style.backgroundColor = "#a0a0a0";
            }
        })
    }
}, 1500);
window.onload =function() {
    $(".contacts_body .contacts li").click(function(e){
        $(".contacts_body .contacts li").removeClass("active");
        $(e.currentTarget).addClass("active");
    });

    $('#action_menu_btn').click(function(){
        $('.action_menu').toggle();
    });
    //搜索框
    $('.contacts_card .search_btn').click(function(){
        console.log(222222334)
    });
    //发送信息
    $('.input-group-append .send_btn').click(function () {
        console.log(222267882222)
        $('.msg_card_body').append('<div class="d-flex justify-content-start mb-4">\n' +
            '                        <div class="img_cont_msg">\n' +
            '                            <img src="https://static.turbosquid.com/Preview/001292/481/WV/_D.jpg" class="rounded-circle user_img_msg">\n' +
            '                       </div>\n' +
            '                        <div class="msg_cotainer">Hi, how are you samim?<span class="msg_time">8:40 AM, Today</span></div></div>')
    });
};
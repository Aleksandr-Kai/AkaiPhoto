$( document ).ready(function() {
    console.log( "ready!" );
    $('#signin-form').on('submit', function(e){
        e.preventDefault(); 
        var public_name = $('#public-name').val();
        var login = $('#inp-login').val();
        var password = $('#inp-pass').val();
        console.log('name=' + public_name);
        console.log('login=' + login);
        console.log('password=' + password);
        if(public_name != '' && login != '' && password != ''){
            $.ajax({
                url:"/sign-in",
                method:"POST",
                data:{name:public_name, password:password, login:login},
                success:function(res){
                    console.log(res)
                    var obj = $.parseJSON(res);
                    
                    if(obj.success){
                        if(obj.message != '') alert(obj.message);
                        $.ajax({
                            url:"/login",
                            method:"POST",
                            data:{password:password, login:login},
                            success:function(res){
                                var obj = $.parseJSON(res);
                                
                                if(obj.success){
                                    if(obj.message != '') alert(obj.message);
                                    location.replace('/home');
                                }
                                else{
                                    if(obj.message != '') alert(obj.message); else alert('Неверный логин или пароль');
                                }
                            }
                        })
                    }
                    else{
                        if(obj.message != '') alert(obj.message); else alert('Неизвестная ошибка');
                    }
                }
            })
        }
        else{
            alert('Нужно заполнить все поля');
        }     
    })
});
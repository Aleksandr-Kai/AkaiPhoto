function onPageLoaded() {
    LoadAlbums(); // загрузка списка альбомов
    //----------------------------------------------------------------------------------------------
    // Событие при отображении окна с контентом альбома
    document.getElementById('myModal').addEventListener('show.bs.modal', function(event) {
        var xhr= new XMLHttpRequest();
        var button = event.relatedTarget
        var url = button.getAttribute('data-href') // здесь ссылка на контент альбома
        xhr.open('GET', url, true);
        xhr.onreadystatechange = function() {
            if (this.readyState!==4) return;
            if (this.status!==200) return;
            document.getElementById('view-album-body').innerHTML = this.responseText;
            document.getElementById('imgs-grid').focus();
            var title = $('h1#gallery-title').text();
            $('h1#modal-title').text(title);    // название альбома в заголовок окна
            $('.gallery-buttons').hide();   // скрыть лишние блоки
            $('img[data-description]').each(function(index){  // задание ссылок для полноэкранного прос
                var src = $(this).attr('src').replace('size=s', 'size=m');
                var obj = $(this).closest('a');
                obj.attr('href', src);
            });
            UpdateFancyBox();
        };
        xhr.send();
    })
    //----------------------------------------------------------------------------------------------
    // Событие при закрытии окна с контентом альбома
    /*
    document.getElementById('myModal').addEventListener('hide.bs.modal', function(event) {
        document.getElementById('mb').innerHTML= "";
    })*/
    //----------------------------------------------------------------------------------------------
    // Событие на кнопку просмотра альбома
    /*
    $('a#card-btn').on('click', function(event){
        $('.album-list').load($(this).attr('data-href'));
    })*/
    //----------------------------------------------------------------------------------------------
    // Обработка кнопки логина
    $('#login-form').on('submit', function(e){
        e.preventDefault(); 
        Login($('#inp-login').val(), $('#inp-pass').val());      
    })
}
//----------------------------------------------------------------------------------------------
function Login(login, password){
    if(login != '' && password != ''){
        $.ajax({
            url:"/login",
            method:"POST",
            data:{password:password, login:login},
            success:function(res){
                //console.log(res)
                var obj = $.parseJSON(res);
                
                if(obj.success){
                    if(obj.message != '') alert(obj.message);
                    location.reload();
                }
                else{
                    if(obj.message != '') alert(obj.message); else alert('Неверный логин или пароль');
                }
            }
        })
    }
    else{
        alert('Нужно заполнить все поля');
    }
}

//----------------------------------------------------------------------------------------------
// Параметры FancyBox
function UpdateFancyBox(){
    $('[data-fancybox="gallery"]').fancybox({
        buttons: [
            "zoom",
            "share",
            "slideShow",
            "fullScreen",
            "download",
            "thumbs",
            "close"
          ],
        afterShow : function(instance, current) {
          var src =  current.src.replace('size=m', 'size=l');
          $("[data-fancybox-download]").attr('href', src);
        },
        loop : true,
        animationEffect: "fade",
    });
}
//----------------------------------------------------------------------------------------------
// Загрузка списка альбомов
function LoadAlbums(){
    $('.album-list').load('home?get_content=album-list', function(){
        $('img[src=""]').each(function(index){
            $(this).parent().parent().hide();
        });
    });
    
}
//----------------------------------------------------------------------------------------------
// Подгрузка окна настроек
function LoadConfig(){
    var m = $('#modals');
    console.log(m);
    m.load('home?get_content=config', function(){
        console.log($('#modal-config'));
        $('#modal-config').modal('show');
    })
}
//----------------------------------------------------------------------------------------------
// Возврат к списку альбомов
function GoBack(){
    LoadAlbums();
    $([document.documentElement, document.body]).animate({
        scrollTop: 0
    }, 200);
}
//----------------------------------------------------------------------------------------------
// Для прокрутки наверх
function ScrollToTop(){
    $([document.documentElement, document.body]).animate({
        scrollTop: 0
    }, 200);
}
//----------------------------------------------------------------------------------------------
// Просмотр альбома в блоке списка альбомов
function OpenGallery(gallery){
    $('.album-list').load(gallery, function(){UpdateFancyBox();});
    $([document.documentElement, document.body]).animate({
        scrollTop: 0
    }, 200);
}
//----------------------------------------------------------------------------------------------
if (document.readyState === 'complete' ||
    (document.readyState !== 'loading' && !document.documentElement.doScroll)) {
    onPageLoaded()
} else {
document.addEventListener('DOMContentLoaded', () => {
    onPageLoaded()
});
}


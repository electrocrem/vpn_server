$('.get-btn').on('click', function() {
    $.ajax({
        url: '/profile',
        method: 'GET',
        success: function(data) {
            var a = document.createElement('a');
            var url = '/profile';
            a.href = url;
            a.download = "1"
            document.body.append(a);
            a.click();
            a.remove();
            window.URL.revokeObjectURL(url);
        }
    });
});
$(document).ready(function () {
    $("#create_tariff").click(
        function () {
            $.ajax({
                url: "/api/tariff",
                type: "POST",
                data: $('#tariff_form').serialize(),
                success: function () {
                    location.replace('/admin/tariffs');
                },
                error: function (xhr, status, error) {
                    var err = eval("(" + xhr.responseText + ")");
                    var message = typeof err.msg !== 'undefined' && err.msg != null && err.msg != '' ? err.msg : 'Incorrect Data'
                    $('#tariff_result').html('<p>'+message+'</p>');
                }
            });
            return false;
        }
    );

    $("a[id^='delete_tariff']").click(
        function () {
            var id = this.id.replace("delete_tariff_","");
            $.ajax({
                url: "/api/tariff/"+id,
                type: "DELETE",
                success: function () {
                    location.reload();
                },
                error: function (xhr, status, error) {
                    var err = eval("(" + xhr.responseText + ")");
                    var message = typeof err.msg !== 'undefined' && err.msg != null && err.msg != '' ? err.msg : 'Cannot delete plan'
                    $('#tariffs_result').html('<p>'+message+'</p>');
                }
            });
            return false;
        }
    );
});
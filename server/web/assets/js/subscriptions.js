$(document).ready(function () {
	$("input[id^='create_subscription']").click(function () {
		var id = this.id.replace("create_subscription_", "");
		$.ajax({
			url: `/api/customer/${id}/subscription`,
			type: "POST",
			data: $("#subscription_form").serialize(),
			success: function () {
				location.replace(`admin/customer/${id}/subscriptions/`);
			},
			error: function (xhr, status, error) {
				var err = eval("(" + xhr.responseText + ")");
				var message =
					typeof err.msg !== "undefined" && err.msg != null && err.msg != "" ? err.msg : "Incorrect Data";
				$("#subscription_result").html("<p>" + message + "</p>");
			},
		});
		return false;
	});

	$("a[id^='delete_subscription']").click(function () {
		var id = this.id.replace("delete_subscription_", "");
		var customerId = this.dataset.customerId;
		$.ajax({
			url: `/api/customer${customerId}/subscription/${id}`,
			type: "DELETE",
			success: function () {
				location.reload();
			},
			error: function (xhr, status, error) {
				var err = eval("(" + xhr.responseText + ")");
				var message =
					typeof err.msg !== "undefined" && err.msg != null && err.msg != "" ? err.msg : "Cannot delete subscription";
				$("#subscriptions_result").html("<p>" + message + "</p>");
			},
		});
		return false;
	});

	$("a[id^='toggle_subscription_']").click(function () {
		var id = this.id.replace("toggle_subscription_", "");
        var customerId = this.dataset.customerId;
		$.ajax({
			url: `/api/customer/${customerId}/subscription/${id}`,
			type: "PATCH",
			data: {
                tariff_id: this.dataset.tariffId,
                stripe_id: this.dataset.stripeId, 
                status: this.dataset.status == "true" ? "" : "on", 
            },
			success: function () {
				location.reload();
			},
			error: function (xhr, status, error) {
				var err = eval("(" + xhr.responseText + ")");
				var message =
					typeof err.msg !== "undefined" && err.msg != null && err.msg != "" ? err.msg : "Cannot modify subscription";
				$("#subscriptions_result").html("<p>" + message + "</p>");
			},
		});
		return false;
	});

    $("a[id^='toggle_license_']").click(function () {
		var id = this.id.replace("toggle_license_", "");
		$.ajax({
			url: `/api/key/${id}`,
			type: "PATCH",
			data: {
                status: this.dataset.status == "true" ? "" : "on", 
            },
			success: function () {
				location.reload();
			},
			error: function (xhr, status, error) {
				var err = eval("(" + xhr.responseText + ")");
				var message =
					typeof err.msg !== "undefined" && err.msg != null && err.msg != "" ? err.msg : "Cannot modify license";
				$("#subscriptions_result").html("<p>" + message + "</p>");
			},
		});
		return false;
	});

	$("a[id^='delete_license_']").click(function () {
		var id = this.id.replace("delete_license_", "");
		var customerId = this.dataset.customerId;
		$.ajax({
			url: `/api/key/${id}`,
			type: "DELETE",
			success: function () {
				location.reload();
			},
			error: function (xhr, status, error) {
				var err = eval("(" + xhr.responseText + ")");
				var message =
					typeof err.msg !== "undefined" && err.msg != null && err.msg != "" ? err.msg : "Cannot delete subscription";
				$("#subscriptions_result").html("<p>" + message + "</p>");
			},
		});
		return false;
	});
});

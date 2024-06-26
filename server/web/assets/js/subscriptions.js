$(document).ready(function () {
	$("input[id^='create_subscription']").click(function () {
		var customerId = this.id.replace("create_subscription_", "");
		$.ajax({
			url: `/api/customer/${customerId}/subscription`,
			type: "POST",
			data: $("#subscription_form").serialize(),
			success: function () {
				location.replace(`/admin/customer/${customerId}/subscriptions/`);
			},
			error: function (xhr, status, error) {
				var err = eval("(" + xhr.responseText + ")");
				var message =
					typeof err.msg !== "undefined" && err.msg != null && err.msg != "" ? err.msg : "Incorrect Data";
				$("#subscription_result").html("<p><b>" + message + "</b></p>");
			},
		});
		return false;
	});

	$("input[id^='update_subscription']").click(function () {
		var id = this.id.replace("update_subscription_", "");
		var customerId = this.dataset.customerId;
		$.ajax({
			url: `/api/customer/${customerId}/subscription/${id}`,
			type: "PATCH",
			data: $("#subscription_form").serialize(),
			success: function () {
				location.replace(`/admin/customer/${customerId}/subscriptions/`);
			},
			error: function (xhr, status, error) {
				var err = eval("(" + xhr.responseText + ")");
				var message =
					typeof err.msg !== "undefined" && err.msg != null && err.msg != "" ? err.msg : "Incorrect Data";
				$("#subscription_result").html("<p><b>" + message + "</b></p>");
			},
		});
		return false;
	});

	$("a[id^='renew_subscription']").click(function () {
		var id = this.id.replace("renew_subscription_", "");
		if (!confirm(`Are you sure you want to renew ${this.title.replace("Renew ", "")} subscription and all its licenses?`)) {
			return
		}
		$.ajax({
			url: `/api/subscription/${id}/renew`,
			type: "PATCH",
			success: function () {
				location.reload();
			},
			error: function (xhr, status, error) {
				var err = eval("(" + xhr.responseText + ")");
				var message =
					typeof err.msg !== "undefined" && err.msg != null && err.msg != "" ? err.msg : "Cannot renew subscription";
				$("#subscriptions_result").html("<p><b>" + message + "</b></p>");
			},
		});
		return false;
	});

	$("a[id^='delete_subscription']").click(function () {
		var id = this.id.replace("delete_subscription_", "");
		var customerId = this.dataset.customerId;
		if (!confirm(`Are you sure you want to delete ${this.title.replace("Delete ", "")} subscription and all its licenses?`)) {
			return
		}
		$.ajax({
			url: `/api/customer/${customerId}/subscription/${id}`,
			type: "DELETE",
			success: function () {
				location.reload();
			},
			error: function (xhr, status, error) {
				var err = eval("(" + xhr.responseText + ")");
				var message =
					typeof err.msg !== "undefined" && err.msg != null && err.msg != "" ? err.msg : "Cannot delete subscription";
				$("#subscriptions_result").html("<p><b>" + message + "</b></p>");
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
				$("#subscriptions_result").html("<p><b>" + message + "</b></p>");
			},
		});
		return false;
	});
});

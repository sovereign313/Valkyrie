        var supportdivwidth;
        var supportdivheight;

        $(function() {
                supportdivwidth = parseInt($("#supportdiv").css('width'), 10);
                supportdivheight = parseInt($("#supportdiv").css('height'), 10);

                $("#supportdiv").hide();
                $("#supportdiv").draggable();

                $("#supportgrid").click(function(evt) {
			offset = $("#supportgrid").offset();
			l = $("#supportgrid").offset().left - (supportdivwidth / 3); 
			t = $("#supportgrid").offset().top - (supportdivheight / 3) 
			$("#supportdiv").css('left', l); 
			$("#supportdiv").css('top', t); 

			$("#supportdiv").children().each(function() {
				$(this).css("visibility", "hidden");
			});
			$("#supportdiv").show("scale", {}, 1000, function() {
				$("#supportdiv").children().each(function() {
					$(this).css("visibility", "visible");
				});

				$("#firstname").focus();

				t = $("#supportdiv").offset().top;
				bottom = t + $("#supportdiv").outerHeight();
				vpt = $(window).scrollTop();
				vpb = vpt + $(window).height();

				if (bottom > vpt && t < vpb) {
					$("#supportdiv").offset({top: t - supportdivheight / 7 });
				}
			});
		});
	});

	function cancel_support() {
		$("#supportdiv").children().each(function() {
			$(this).css("visibility", "hidden");
		});

		$("#supportdiv").hide("scale", {}, 1000);
	}

	function contact_support() {
		license = localStorage.getItem("licensekey");
		business = localStorage.getItem("business");

		fname = $("#firstname").val();
		lname = $("#lastname").val();
		phone = $("#phonenumber").val();
		email = $("#emailaddress").val();
		company = $("#company").val();
		issue = $("#issue").val();

		if (fname == "") {
                        $().toastmessage('showErrorToast', "First Name Missing");
                        $("#firstname").effect("shake");
                        $("#firstname").focus();
                        return;
		}

		if (lname == "") {
                        $().toastmessage('showErrorToast', "Last Name Missing");
                        $("#lastname").effect("shake");
                        $("#lastname").focus();
                        return;
		}

		if (phone == "") {
                        $().toastmessage('showErrorToast', "Phone Number Missing");
                        $("#phonenumber").effect("shake");
                        $("#phonenumber").focus();
                        return;
		}

		if (email == "") {
                        $().toastmessage('showErrorToast', "Email Address Missing");
                        $("#email").effect("shake");
                        $("#email").focus();
                        return;
		}

		if (company == "") {
                        $().toastmessage('showErrorToast', "Company Is Missing");
                        $("#company").effect("shake");
                        $("#company").focus();
                        return;
		}

		if (issue == "") {
                        $().toastmessage('showErrorToast', "Missing Issue");
                        $("#issue").effect("shake");
                        $("#issue").focus();
                        return;
		}

		var data = {
			fname:fname,
			lname:lname,
			phone:phone,
			email:email,
			company:company,
			issue:issue,
			license:license,
			business:business
		}

		$.post(
			'/contact',
			data,
			function(responseText) {
				resp = JSON.parse(responseText);
				if (resp.Message == "success") {
					cancel_support();
					$().toastmessage('showSuccessToast', "Notified Support!  Someone Will Be With You Shortly!");
				} else if (resp.Code == "502") {
					$().toastmessage('showErrorToast', "Tampering Has Been Detected.  Please Contact Support");
				} else {
					$().toastmessage('showErrorToast', 'There Was An Error: ' + resp.Message);
					$().toastmessage('showErrorToast', 'Please Call Us, And Let Us Know');
				}
			}
		);
	}

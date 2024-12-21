	class Alerter {
		constructor(emailserver, fromaddress, twilioaccount, twiliotoken, twiliophonenumber, host) {
			this.EmailServer = emailserver;
			this.FromAddress = fromaddress;
			this.TwilioAccount = twilioaccount;
			this.TwilioToken = twiliotoken;
			this.TwilioPhoneNumber = twiliophonenumber
			this.Host = host;
		}

		Complete() {
			if (awsconfig != null) {
				return true;
			}

			if (this.EmailServer == "" && this.TwilioAccount == "") {
				return false;
			}

			if (this.EmailServer != "" && this.FromAddress == "") {
				return false;
			}

			if (this.TwilioAccount != "" && this.TwilioToken == "") {
				return false;
			}

			if (this.TwilioAccount != "" && this.TwilioPhoneNumber == "") {
				return false;
			}

			if (this.Host == "" || this.Host == "none") {
				return false;
			}
			
			return true;
		}
	}

        $(function() {

	        alerterconfig = JSON.parse(localStorage.getItem("alerterconfig"));
	        awsconfig = JSON.parse(localStorage.getItem("awsconfig"));
		alerterdivwidth = parseInt($("#alerterdiv").css('width'), 10);
                alerterdivheight = parseInt($("#alerterdiv").css('height'), 10);

		if (alerterconfig != null) {
			$("#alertergrid").removeClass('grid-item-goodenough');
			$("#alertergrid").addClass('grid-item-complete');
		}

		if (awsconfig != null) {
			$('#awsconfigureddiv').css('display', 'block');
		}

		if (awsconfig != null && alerterconfig == null) {
			$('#alertergrid').tooltip('hide').attr('title', 'Configured To Use SNS. Further Configuration Is Not Required.').tooltip('fixTitle');
			$('#alertergrid').addClass('grid-item-goodenough');
		}


		$("#alerterdiv").hide();
		$("#alerterdiv").draggable();
	
		$("#alertergrid").click(function(evt) {
	                offset = $("#alertergrid").offset();
	                l = $("#alertergrid").offset().left - (alerterdivwidth / 3);
        	        t = $("#alertergrid").offset().top - (alerterdivheight / 3) 
                	$("#alerterdiv").css('left', l);
	                $("#alerterdiv").css('top', t);

                        hosts = JSON.parse(localStorage.getItem("hostsconfig"));
			$("#alerterhost").empty();
                        if (hosts != null) {
                                for (i = 0; i < hosts.length; i++) {
                                        option = document.createElement("option");
                                        option.text = hosts[i];
                                        option.value = hosts[i];
                                        document.getElementById('alerterhost').add(option);
                                }
                        } else {
                                option = document.createElement("option");
                                option.text = "Hosts Not Configured";
                                option.value = "none";
                                document.getElementById('alerterhost').add(option);
                        }


			if (alerterconfig != null) {
				$("#emailserver").val(alerterconfig.EmailServer);
				$("#fromaddress").val(alerterconfig.FromAddress);
				$("#twilioaccount").val(alerterconfig.TwilioAccount);
				$("#twiliotoken").val(alerterconfig.TwilioToken);
				$("#twiliophonenumber").val(alerterconfig.TwilioPhoneNumber);
				$("#alerterhost").val(alerterconfig.Host);
			}
	
			$("#alerterdiv").children().each(function() {
				$(this).css("visibility", "hidden");
			});

			$("#alerterdiv").show("scale", {}, 1000, function() {
				$("#alerterdiv").children().each(function() {
					$(this).css("visibility", "visible");
				});

				$("#alerterregion").focus();
			});
		});

		$("#emailserver").keypress(function(event) {
			var keycode = (event.keyCode ? event.keyCode : event.which);
			if (keycode == '13') {
				save_alerter();
				return;
			}
		});

		$("#fromaddress").keypress(function(event) {
			var keycode = (event.keyCode ? event.keyCode : event.which);
			if (keycode == '13') {
				save_alerter();
				return;
			}
		});

		$("#twilioaccount").keypress(function(event) {
			var keycode = (event.keyCode ? event.keyCode : event.which);
			if (keycode == '13') {
				save_alerter();
				return;
			}
		});

		$("#twilioaccount").keypress(function(event) {
			var keycode = (event.keyCode ? event.keyCode : event.which);
			if (keycode == '13') {
				save_alerter();
				return;
			}
		});

		$("#twiliotoken").keypress(function(event) {
			var keycode = (event.keyCode ? event.keyCode : event.which);
			if (keycode == '13') {
				save_alerter();
				return;
			}
		});

		$("#twiliophonenumber").keypress(function(event) {
			var keycode = (event.keyCode ? event.keyCode : event.which);
			if (keycode == '13') {
				save_alerter();
				return;
			}
		});
	});

        function cancel_alerter() {
                $("#alerterdiv").children().each(function() {
                        $(this).css("visibility", "hidden");
                });

		$("#alerterdiv").find('input[type="text"]').val('');
                $("#alerterdiv").hide("scale", {}, 1000);
        }

        function save_alerter() {
		emailserver = $("#emailserver").val();
		fromaddress = $("#fromaddress").val();
		twilioaccount = $("#twilioaccount").val();
		twiliotoken = $("#twiliotoken").val();
		twiliophonenumber = $("#twiliophonenumber").val();
		alerterhost = $("#alerterhost").val();

		if (emailserver == "" && twilioaccount == "") {
			$().toastmessage('showErrorToast', "You Must Configure Either Email or Twilio");
			$().toastmessage('showErrorToast', "If you configured AWS, then just cancel and move on");
                        $("#emailserver").effect("shake");
                        $("#twilioaccount").effect("shake");
                        return;
		}

		if (emailserver != "" && fromaddress == "") {
			$().toastmessage('showErrorToast', "Missing From Address");
                        $("#fromaddress").effect("shake");
                        return;
		}

		if (emailserver != "" && ! emailserver.includes(":")) {
                        $().toastmessage('showErrorToast', "You Must Include The Port Number: host:port");
                        $("#emailserver").effect("shake");
                        return;
		}

		if (twilioaccount != "" && twiliotoken == "") {
			$().toastmessage('showErrorToast', "Missing Twilio Token");
                        $("#twiliotoken").effect("shake");
                        return;
		}

		if (twilioaccount != "" && twiliophonenumber == "") {
			$().toastmessage('showErrorToast', "Missing Your Twilio Phone Number");
                        $("#twiliophonenumber").effect("shake");
                        return;
		}

		if (alerterhost == "none") {
                        $().toastmessage('showErrorToast', "Missing Host To Install On");
                        $("#alerterhost").effect("shake");
                        return;
		}

		alerterconfig = new Alerter(emailserver, fromaddress, twilioaccount, twiliotoken, twiliophonenumber, alerterhost);
		if (! alerterconfig.Complete()) {
			$().toastmessage('showErrorToast', "Error Saving Alerter Data!!!");
                        $("#alerterdiv").effect("shake");
                        return;
		}

		localStorage.setItem("alerterconfig", JSON.stringify(alerterconfig))
		$("#alertergrid").removeClass('grid-item-goodenough');
		$("#alertergrid").addClass('grid-item-complete');

                cancel_alerter();
                $().toastmessage('showSuccessToast', "Saved Alerter Information");

		if (isInstallReady()) {
			$('#install').tooltip('hide').attr('title', 'Install Valkyrie!').tooltip('fixTitle');
			$('#generatemanifest').tooltip('hide').attr('title', 'Generate An Installation Manifest!').tooltip('fixTitle');
		}

        }




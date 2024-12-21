	class MailReader {
		constructor(dupprotection, usetls, mailprotocol, protecttime, sleeptimeout, mailserver, mailuser, mailpassword, mailsubjecttrigger, host) {
			this.DupProtection = dupprotection;
			this.UseTLS = usetls;
			this.MailProtocol = mailprotocol;
			this.ProtectTime = protecttime;
			this.SleepTimeout = sleeptimeout;
			this.MailServer = mailserver;
			this.MailUser = mailuser;
			this.MailPassword = mailpassword;
			this.MailSubjectTrigger = mailsubjecttrigger;
			this.Host = host;
		}

		Complete() {
			if (this.DupProtection == "true" && this.ProtectTime == "") {
				return false;
			}

			if (this.SleepTimeout == "") {
				return false;
			}

			if (this.MailServer == "") {
				return false;
			}

			if (this.MailUser == "") {
				return false;
			}

			if (this.MailPassword == "") {
				return false;
			}

			if (this.MailSubjectTrigger == "") {
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

	        mailreaderconfig = JSON.parse(localStorage.getItem("mailreaderconfig"));
		mailreaderdivwidth = parseInt($("#mailreaderdiv").css('width'), 10);
                mailreaderdivheight = parseInt($("#mailreaderdiv").css('height'), 10);

		if (mailreaderconfig != null) {
			$("#mailreadergrid").addClass('grid-item-complete');
		}

		$("#mailreaderdiv").hide();
		$("#mailreaderdiv").draggable();

                $("#duplicationprotectionchk").change(function(evt) {
                        if (this.checked) {
                                $("#protecttime").removeAttr('disabled');
                                $("#protecttime").focus();
                        } else {
                                $("#protecttime").attr('disabled', 'disabled');
                        }
                });
	
		$("#mailreadergrid").click(function(evt) {
	                offset = $("#mailreadergrid").offset();
	                l = $("#mailreadergrid").offset().left - (mailreaderdivwidth / 3);
        	        t = $("#mailreadergrid").offset().top - (mailreaderdivheight / 3) 
                	$("#mailreaderdiv").css('left', l);
	                $("#mailreaderdiv").css('top', t);

                        hosts = JSON.parse(localStorage.getItem("hostsconfig"));
                        $("#mailreaderhost").empty();

                        if (hosts != null) {
                                for (i = 0; i < hosts.length; i++) {
                                        option = document.createElement("option");
                                        option.text = hosts[i];
                                        option.value = hosts[i];
                                        document.getElementById('mailreaderhost').add(option);
                                }
                        } else {
                                option = document.createElement("option");
                                option.text = "Hosts Not Configured";
                                option.value = "none";
                                document.getElementById('mailreaderhost').add(option);
                        }

			if (mailreaderconfig != null) {
				if (mailreaderconfig.DupProtection == "true") {
                                        $("#duplicationprotectionchk").prop('checked', true);
                                        $("#protecttime").removeAttr('disabled');
                                        $("#protecttime").val(mailreaderconfig.ProtectTime);
				} else {
                                        $("#duplicationprotectionchk").prop('checked', false);
                                        if (mailreaderconfig.ProtectTime != "") {
                                                $("#protecttime").val(mailreaderconfig.ProtectTime);
                                        }
                                        $("#protecttime").attr('disabled', 'disabled');
                                }

				if (mailreaderconfig.UseTLS == "true") {
					$("#usetlschk").prop('checked', true);
				} else {
					$("#usetlschk").prop('checked', true);
				}

				if (mailreaderconfig.MailProtocol == "imap") {
					$("#mailprotocol").val('imap');
				} else {
					$("#mailprotocol").val('pop');
				}

				$("#protecttime").val(mailreaderconfig.ProtectTime);
				$("#sleeptimeout").val(mailreaderconfig.SleepTimeout);
				$("#mailserver").val(mailreaderconfig.MailServer);
				$("#mailuser").val(mailreaderconfig.MailUser);
				$("#mailpassword").val(mailreaderconfig.MailPassword);
				$("#mailsubjecttrigger").val(mailreaderconfig.MailSubjectTrigger);
                                $("#mailreaderhost").val(mailreaderconfig.Host);
			}
	
			$("#mailreaderdiv").children().each(function() {
				$(this).css("visibility", "hidden");
			});

			$("#mailreaderdiv").show("scale", {}, 1000, function() {
				$("#mailreaderdiv").children().each(function() {
					$(this).css("visibility", "visible");
				});
			});
		});

	});

        function cancel_mailreader() {
                $("#mailreaderdiv").children().each(function() {
                        $(this).css("visibility", "hidden");
                });

                $("#duplicationprotectionchk").prop('checked', false);
                $("#usetlschk").prop('checked', false);
		$("#mailprotocol").val('imap');
                $("#mailreaderdiv").find('input[type="text"]').val('');
                $("#mailreaderdiv").hide("scale", {}, 1000);
        }

        function save_mailreader() {
                dupprotection = $("#duplicationprotectionchk").prop('checked') ? "true" : "false";
                usetls = $("#usetlschk").prop('checked') ? "true" : "false";
		mailprotocol = $("#mailprotocol").val();
		protecttime = $("#protecttime").val();
		sleeptimeout = $("#sleeptimeout").val();
		mailserver = $("#mailserver").val();
		mailuser = $("#mailuser").val();
		mailpassword = $("#mailpassword").val();
		mailsubjecttrigger = $("#mailsubjecttrigger").val();
                mailreaderhost = $("#mailreaderhost").val();

		if (dupprotection == "true" && protecttime == "") {
			$().toastmessage('showErrorToast', "Missing Protection Time");
                        $("#protecttime").effect("shake");
                        return;
		}

		if (sleeptimeout == "") {
			$().toastmessage('showErrorToast', "Missing Sleep Timeout");
                        $("#sleeptimeout").effect("shake");
                        return;
		}

		if (mailserver == "") {
			$().toastmessage('showErrorToast', "Missing Mail Server");
                        $("#mailserver").effect("shake");
                        return;
		}

		if (mailuser == "") {
			$().toastmessage('showErrorToast', "Missing Mail User");
                        $("#mailuser").effect("shake");
                        return;
		}

		if (mailpassword == "") {
			$().toastmessage('showErrorToast', "Missing Mail Password");
                        $("#mailpassword").effect("shake");
                        return;
		}

		if (mailsubjecttrigger == "") {
			$().toastmessage('showErrorToast', "Missing Mail Subject Trigger");
                        $("#mailsubjecttrigger").effect("shake");
                        return;
		}

                if (mailreaderhost == "none") {
                        $().toastmessage('showErrorToast', "Missing Host To Install On");
                        $("#mailreaderhost").effect("shake");
                        return;
                }

		if (! mailserver.includes(":")) {
                        $().toastmessage('showErrorToast', "You Must Include The Port Number: host:port");
                        $("#mailserver").effect("shake");
                        return;
		}


		mailreaderconfig = new MailReader(dupprotection, usetls, mailprotocol, protecttime, sleeptimeout, mailserver, mailuser, mailpassword, mailsubjecttrigger, mailreaderhost);
		if (! mailreaderconfig.Complete()) {
			$().toastmessage('showErrorToast', "Error Saving MailReader Data!!!");
                        $("#mailreaderdiv").effect("shake");
                        return;
		}

		localStorage.setItem("mailreaderconfig", JSON.stringify(mailreaderconfig))
		$("#mailreadergrid").addClass('grid-item-complete');
                cancel_mailreader();
                $().toastmessage('showSuccessToast', "Saved MailReader Information");
        
		if (isInstallReady()) {
                        $('#install').tooltip('hide').attr('title', 'Install Valkyrie!').tooltip('fixTitle');
                        $('#generatemanifest').tooltip('hide').attr('title', 'Generate An Installation Manifest!').tooltip('fixTitle');
                }
        }




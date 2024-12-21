	class Logger {
		constructor(useeventstreams, usesecurelogging, eshostport, logkey, logfilelocation, host) {
			this.UseEventStreams = useeventstreams;
			this.UseSecureLogging = usesecurelogging;
			this.ESHostPort = eshostport;
			this.LogKey = logkey;
			this.LogFileLocation = logfilelocation;
			this.Host = host;
		}

		Complete() {
			if (this.UseEventStreams == "true" && this.ESHostPort == "") {
				return false;
			}

			if (this.UseSecureLogging == "true" && this.LogKey == "") {
				return false;
			}

			if (this.LogFileLocation == "") {
				return false;
			}

			if (this.Host == "" || this.Host == "none") {
				return false;
			}

			return true;
		}
	}

        $(function() {
	        loggerconfig = JSON.parse(localStorage.getItem("loggerconfig"));
		loggerdivwidth = parseInt($("#loggerdiv").css('width'), 10);
                loggerdivheight = parseInt($("#loggerdiv").css('height'), 10);

		if (loggerconfig != null) {
			$("#loggergrid").addClass('grid-item-complete');
		}

		$("#loggerdiv").hide();
		$("#loggerdiv").draggable();

		$("#eventstreamschk").change(function(evt) {
			if (this.checked) {
				$("#eshostport").removeAttr('disabled');
				$("#eshostport").focus();
			} else {
				$("#eshostport").attr('disabled', 'disabled');
			}
		});
	
		$("#secureloggingchk").change(function(evt) {
			if (this.checked) {
				$("#logkey").removeAttr('disabled');
				$("#logkey").focus();
			} else {
				$("#logkey").attr('disabled', 'disabled');
			}
		});

		$("#loggergrid").click(function(evt) {
	                offset = $("#loggergrid").offset();
	                l = $("#loggergrid").offset().left - (loggerdivwidth / 3);
        	        t = $("#loggergrid").offset().top - (loggerdivheight / 3) 
                	$("#loggerdiv").css('left', l);
	                $("#loggerdiv").css('top', t);

                        hosts = JSON.parse(localStorage.getItem("hostsconfig"));
			$("#loggerhost").empty();
                        if (hosts != null) {
                                for (i = 0; i < hosts.length; i++) {
                                        option = document.createElement("option");
                                        option.text = hosts[i];
                                        option.value = hosts[i];
                                        document.getElementById('loggerhost').add(option);
                                }
                        } else {
                                option = document.createElement("option");
                                option.text = "Hosts Not Configured";
                                option.value = "none";
                                document.getElementById('loggerhost').add(option);
                        }

			if (loggerconfig != null) {
				if (loggerconfig.UseEventStreams == "true") {
					$("#eventstreamschk").prop('checked', true);
					$("#eshostport").removeAttr('disabled');
					$("#eshostport").val(loggerconfig.ESHostPort);
				} else {
					$("#eventstreamschk").prop('checked', false);
					if (loggerconfig.ESHostPort != "") {
						$("#eshostport").val(loggerconfig.ESHostPort);
					}
					$("#eshostport").attr('disabled', 'disabled');
				}

				if (loggerconfig.UseSecureLogging == "true") {
					$("#secureloggingchk").prop('checked', true);
					$("#logkey").removeAttr('disabled');
					$("#logkey").val(loggerconfig.LogKey);
				} else {
					$("#secureloggingchk").prop('checked', false);
					if (loggerconfig.ESHostPort != "") {
						$("#logkey").val(loggerconfig.LogKey);
					}
					$("#logkey").attr('disabled', 'disabled');
				}

				$("#logfilelocation").val(loggerconfig.LogFileLocation);
                                $("#loggerhost").val(loggerconfig.Host);
			}
	
			$("#loggerdiv").children().each(function() {
				$(this).css("visibility", "hidden");
			});

			$("#loggerdiv").show("scale", {}, 1000, function() {
				$("#loggerdiv").children().each(function() {
					$(this).css("visibility", "visible");
				});

				$("#logfilelocation").focus();
			});
		});
	});

        function cancel_logger() {
                $("#loggerdiv").children().each(function() {
                        $(this).css("visibility", "hidden");
                });

                $("#loggerdiv").hide("scale", {}, 1000);
        }

        function save_logger() {
		evtstreams = $("#eventstreamschk").prop('checked') ? "true" : "false";
		securelogging = $("#secureloggingchk").prop('checked') ? "true" : "false";
		eshostport = $("#eshostport").val();
		logkey = $("#logkey").val();
		logfilelocation = $("#logfilelocation").val();
                loggerhost = $("#loggerhost").val();

		if (evtstreams == "true" && eshostport == "") {
			$().toastmessage('showErrorToast', "Missing Host & Port");
                        $("#eshostport").effect("shake");
                        return;
		}

		if (evtstreams == "true" && ! eshostport.includes(":")) {
			$().toastmessage('showErrorToast', "You Must Include The Port Number: host:port");
			$("#eshostport").effect("shake");
			return;
		}

		if (securelogging == "true" && logkey == "") {
			$().toastmessage('showErrorToast', "Missing Secure Log");
                        $("#logkey").effect("shake");
                        return;
		}

		if (logfilelocation == "") {
			$().toastmessage('showErrorToast', "Missing Log File Location");
                        $("#logfilelocation").effect("shake");
                        return;
		}

                if (loggerhost == "none") {
                        $().toastmessage('showErrorToast', "Missing Host To Install On");
                        $("#loggerhost").effect("shake");
                        return;
                }

		loggerconfig = new Logger(evtstreams, securelogging, eshostport, logkey, logfilelocation, loggerhost);
		if (! loggerconfig.Complete()) {
			$().toastmessage('showErrorToast', "Error Saving Logger Data!!!");
                        $("#loggerdiv").effect("shake");
                        return;
		}

		localStorage.setItem("loggerconfig", JSON.stringify(loggerconfig))
		$("#loggergrid").addClass('grid-item-complete');
                cancel_logger();
                $().toastmessage('showSuccessToast', "Saved Logger Information");
        }




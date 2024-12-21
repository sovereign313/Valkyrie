       class Foreman {
                constructor(dupprotection, protecttime, host) { 
                        this.DupProtection = dupprotection;
                        this.ProtectTime = protecttime;
			this.Host = host
                }

                Complete() {
                        if (this.DupProtection == "true" && this.ProtectTime == "" && this.Host == "") {
                                return false;
                        }

                        return true;
                }
        }

        var foremandivwidth;
        var foremandivheight;

        $(function() {

		foremanconfig = JSON.parse(localStorage.getItem("foremanconfig"));
                foremandivwidth = parseInt($("#foremandiv").css('width'), 10);
                foremandivheight = parseInt($("#foremandiv").css('height'), 10);

		if (foremanconfig != null) {
	                $("#foremangrid").addClass('grid-item-complete');
		}

                $("#foremandiv").hide();
                $("#foremandiv").draggable();

                $("#foremandupchk").change(function(evt) {
                        if (this.checked) {
                                $("#foremanprotecttime").removeAttr('disabled');
                                $("#foremanprotecttime").focus();
                        } else {
                                $("#foremanprotecttime").attr('disabled', 'disabled');
                        }
                });

                $("#foremangrid").click(function(evt) {
			offset = $("#foremangrid").offset();
			l = $("#foremangrid").offset().left - (foremandivwidth / 3); 
			t = $("#foremangrid").offset().top - (foremandivheight / 3) 
			$("#foremandiv").css('left', l); 
			$("#foremandiv").css('top', t); 

			hosts = JSON.parse(localStorage.getItem("hostsconfig"));
			$("#foremanhost").empty();

			if (hosts != null) {
				for (i = 0; i < hosts.length; i++) {
					option = document.createElement("option");
			                option.text = hosts[i];
					option.value = hosts[i];
			                document.getElementById('foremanhost').add(option);
				}
			} else {
				option = document.createElement("option");
				option.text = "Hosts Not Configured"; 
				option.value = "none";
				document.getElementById('foremanhost').add(option);
			}

			if (foremanconfig != null) {
                                if (foremanconfig.DupProtection == "true") {
                                        $("#foremandupchk").prop('checked', true);
                                        $("#foremanprotecttime").removeAttr('disabled');
                                        $("#foremanprotecttime").val(foremanconfig.ProtectTime);
                                } else {
                                        $("#foremandupchk").prop('checked', false);
                                        if (foremanconfig.ProtectTime != "") {
                                                $("#foremanprotecttime").val(foremanconfig.ProtectTime);
                                        }
                                        $("#foremanprotecttime").attr('disabled', 'disabled');
                                }

				$("#foremanhost").val(foremanconfig.Host);
			}

			$("#foremandiv").children().each(function() {
				$(this).css("visibility", "hidden");
			});

			$("#foremandiv").show("scale", {}, 1000, function() {
				$("#foremandiv").children().each(function() {
					$(this).css("visibility", "visible");
				});
			});
		});
	});

	function cancel_foreman() {
		$("#foremandiv").children().each(function() {
			$(this).css("visibility", "hidden");
		});

                $("#foremandupchk").prop('checked', false);
		$("#foremandiv").hide("scale", {}, 1000);
		$('#foremanprotecttime').val('');
	}

	function save_foreman() {
                dupprotection = $("#foremandupchk").prop('checked') ? "true" : "false";
		protecttime = $("#foremanprotecttime").val();
		foremanhost = $("#foremanhost").val();

		if (dupprotection == "true" && protecttime == "") {
			$().toastmessage('showErrorToast', "Missing Protection Time");
                        $("#foremanprotecttime").effect("shake");
                        return;
		}

		if (foremanhost == "none") {
			$().toastmessage('showErrorToast', "Missing Host To Install On");
			$("#foremanhost").effect("shake");
			return;
		}

		foremanconfig = new Foreman(dupprotection, protecttime, foremanhost);
		if (! foremanconfig.Complete()) {
			$().toastmessage('showErrorToast', "Error Saving Foreman Data!!!");
                        $("#foremandiv").effect("shake");
                        return;
		}

		localStorage.setItem("foremanconfig", JSON.stringify(foremanconfig))
		$("#foremangrid").addClass('grid-item-complete');
                cancel_foreman();
                $().toastmessage('showSuccessToast', "Saved Foreman Information");
        
	        if (isInstallReady()) {
                        $('#install').tooltip('hide').attr('title', 'Install Valkyrie!').tooltip('fixTitle');
                        $('#generatemanifest').tooltip('hide').attr('title', 'Generate An Installation Manifest!').tooltip('fixTitle');
                }
	}

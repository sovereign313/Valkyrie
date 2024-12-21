       class Dispatcher {
                constructor(dupprotection, protecttime, host) { 
                        this.DupProtection = dupprotection;
                        this.ProtectTime = protecttime;
			this.Host = host;
                }

                Complete() {
                        if (this.DupProtection == "true" && this.ProtectTime == "" && this.Host == "") {
                                return false;
                        }

                        return true;
                }
        }

        var dispatcherdivwidth;
        var dispatcherdivheight;

        $(function() {

		dispatcherconfig = JSON.parse(localStorage.getItem("dispatcherconfig"));
                awsconfig = JSON.parse(localStorage.getItem("awsconfig"));

                dispatcherdivwidth = parseInt($("#dispatcherdiv").css('width'), 10);
                dispatcherdivheight = parseInt($("#dispatcherdiv").css('height'), 10);

		if (awsconfig != null) {
			$('#awsconfigureddispatcherdiv').css('display', 'block');
		} else {
			$('#awsconfigureddispatcherdiv').css('color', '#ff0000');
			$('#awsconfigureddispatcherdiv').text("AWS IS NOT CONFIGURED!");
			$('#awsconfigureddispatcherdiv').css('display', 'block');
		}

		if (dispatcherconfig != null) {
			if (awsconfig != null) {
		                $("#dispatchergrid").addClass('grid-item-complete');
			}
		} 

                $("#dispatcherdiv").hide();
                $("#dispatcherdiv").draggable();

                $("#dispatcherdupchk").change(function(evt) {
                        if (this.checked) {
                                $("#dispatcherprotecttime").removeAttr('disabled');
                                $("#dispatcherprotecttime").focus();
                        } else {
                                $("#dispatcherprotecttime").attr('disabled', 'disabled');
                        }
                });

                $("#dispatchergrid").click(function(evt) {
			offset = $("#dispatchergrid").offset();
			l = $("#dispatchergrid").offset().left - (dispatcherdivwidth / 3); 
			t = $("#dispatchergrid").offset().top - (dispatcherdivheight / 3) 
			$("#dispatcherdiv").css('left', l); 
			$("#dispatcherdiv").css('top', t); 

                	awsconfig = JSON.parse(localStorage.getItem("awsconfig"));
			if (awsconfig != null) {
				$('#awsconfigureddispatcherdiv').css('display', 'block');
				$('#awsconfigureddispatcherdiv').css('color', '#ffff00');
				$('#awsconfigureddispatcherdiv').text("AWS Is Already Configured!");
			} else {
				$('#awsconfigureddispatcherdiv').css('color', '#ff0000');
				$('#awsconfigureddispatcherdiv').text("AWS IS NOT CONFIGURED!");
				$('#awsconfigureddispatcherdiv').css('display', 'block');
			}

                        hosts = JSON.parse(localStorage.getItem("hostsconfig"));
                        $("#dispatcherhost").empty();

                        if (hosts != null) {
                                for (i = 0; i < hosts.length; i++) {
                                        option = document.createElement("option");
                                        option.text = hosts[i];
                                        option.value = hosts[i];
                                        document.getElementById('dispatcherhost').add(option);
                                }
                        } else {
                                option = document.createElement("option");
                                option.text = "Hosts Not Configured";
                                option.value = "none";
                                document.getElementById('dispatcherhost').add(option);
                        }

			if (dispatcherconfig != null) {
                                if (dispatcherconfig.DupProtection == "true") {
                                        $("#dispatcherdupchk").prop('checked', true);
                                        $("#dispatcherprotecttime").removeAttr('disabled');
                                        $("#dispatcherprotecttime").val(dispatcherconfig.ProtectTime);
                                } else {
                                        $("#dispatcherdupchk").prop('checked', false);
                                        if (dispatcherconfig.ProtectTime != "") {
                                                $("#dispatcherprotecttime").val(dispatcherconfig.ProtectTime);
                                        }
                                        $("#dispatcherprotecttime").attr('disabled', 'disabled');
                                }

                                $("#dispatcherhost").val(dispatcherconfig.Host);
			}

			$("#dispatcherdiv").children().each(function() {
				$(this).css("visibility", "hidden");
			});

			$("#dispatcherdiv").show("scale", {}, 1000, function() {
				$("#dispatcherdiv").children().each(function() {
					$(this).css("visibility", "visible");
				});
			});
		});
	});

	function cancel_dispatcher() {
		$("#dispatcherdiv").children().each(function() {
			$(this).css("visibility", "hidden");
		});

                $("#dispatcherdupchk").prop('checked', false);
		$("#dispatcherdiv").hide("scale", {}, 1000);
		$('#dispatcherprotecttime').val('');
	}

	function save_dispatcher() {
                dupprotection = $("#dispatcherdupchk").prop('checked') ? "true" : "false";
		protecttime = $("#dispatcherprotecttime").val();
                dispatcherhost = $("#dispatcherhost").val();

		if (dupprotection == "true" && protecttime == "") {
			$().toastmessage('showErrorToast', "Missing Protection Time");
                        $("#dispatcherprotecttime").effect("shake");
                        return;
		}

                if (dispatcherhost == "none") {
                        $().toastmessage('showErrorToast', "Missing Host To Install On");
                        $("#dispatcherhost").effect("shake");
                        return;
                }

		dispatcherconfig = new Dispatcher(dupprotection, protecttime, dispatcherhost);
		if (! dispatcherconfig.Complete()) {
			$().toastmessage('showErrorToast', "Error Saving Dispatcher Data!!!");
                        $("#dispatcherdiv").effect("shake");
                        return;
		}

		localStorage.setItem("dispatcherconfig", JSON.stringify(dispatcherconfig))
		if (awsconfig != null) {
			$("#dispatchergrid").addClass('grid-item-complete');
		}
                cancel_dispatcher();
                $().toastmessage('showSuccessToast', "Saved Dispatcher Information");

                if (isInstallReady()) {
                        $('#install').tooltip('hide').attr('title', 'Install Valkyrie!').tooltip('fixTitle');
                        $('#generatemanifest').tooltip('hide').attr('title', 'Generate An Installation Manifest!').tooltip('fixTitle');
                }
	}

       class Launcher {
                constructor(launcherdefaultimage, host) { 
                        this.DefaultImage = launcherdefaultimage;
			this.Host = [];
			
			for (i=0; i < host.length; i++) {
				this.Host.push(host[i])
			}
                }

                Complete() {
                        if (this.DefaultImage == "" || this.Host == "") {
                                return false;
                        }

                        return true;
                }
        }

        var launcherdivwidth;
        var launcherdivheight;

        $(function() {

		launcherconfig = JSON.parse(localStorage.getItem("launcherconfig"));
                launcherdivwidth = parseInt($("#launcherdiv").css('width'), 10);
                launcherdivheight = parseInt($("#launcherdiv").css('height'), 10);

		if (launcherconfig != null) {
	                $("#launchergrid").addClass('grid-item-complete');
		}

                $("#launcherdiv").hide();
                $("#launcherdiv").draggable();

		$("#launcherhost").click(function() {
			hst = $("#launcherhost").val();
			if (hst == null) return;
			$('option:selected', this).remove();
			option = document.createElement("option");
			option.text = hst; 
			option.value = hst;
			document.getElementById('launcherworkers').add(option);
		});

		$("#launcherworkers").click(function() {
			hst = $("#launcherworkers").val();
			if (hst == null) return;
			$('option:selected', this).remove();
			option = document.createElement("option");
			option.text = hst; 
			option.value = hst;
			document.getElementById('launcherhost').add(option);
		});

                $("#launchergrid").click(function(evt) {
			offset = $("#launchergrid").offset();
			l = $("#launchergrid").offset().left - (launcherdivwidth / 3); 
			t = $("#launchergrid").offset().top - (launcherdivheight / 3) 
			$("#launcherdiv").css('left', l); 
			$("#launcherdiv").css('top', t); 

			$("#launcherhost").empty();
			$("#launcherworkers").empty();

                        hosts = JSON.parse(localStorage.getItem("hostsconfig"));
			$("#launcherhost").empty();

                        if (hosts != null) {
                                for (i = 0; i < hosts.length; i++) {
                                        option = document.createElement("option");
                                        option.text = hosts[i];
                                        option.value = hosts[i];
                                        document.getElementById('launcherhost').add(option);
                                }
                        } else {
				$().toastmessage('showErrorToast', "Hosts Are Not Configured");
                        }


			if (launcherconfig != null) {
				$("#launcherdefaultimage").val(launcherconfig.DefaultImage);
				for (i=0; i < launcherconfig.Host.length; i++) {
					option = document.createElement("option");
                                        option.text = hosts[i];
                                        option.value = hosts[i];
                                        document.getElementById('launcherworkers').add(option);
					$("#launcherhost option[value='" + hosts[i] + "']").remove();
				}
			}

			$("#launcherdiv").children().each(function() {
				$(this).css("visibility", "hidden");
			});

			$("#launcherdiv").show("scale", {}, 1000, function() {
				$("#launcherdiv").children().each(function() {
					$(this).css("visibility", "visible");
				});
			});
		});
	});

	function cancel_launcher() {
		$("#launcherdiv").children().each(function() {
			$(this).css("visibility", "hidden");
		});

		$("#launcherdiv").hide("scale", {}, 1000);
		$('#launcherdefaultimage').val('');
	}

	function save_launcher() {
		launcherdefaultimage = $("#launcherdefaultimage").val();

		if (launcherdefaultimage == "") {
			$().toastmessage('showErrorToast', "Missing Default Image");
                        $("#launcherdefaultimage").effect("shake");
                        return;
		}

                lhosts = [];
                $('#launcherworkers').find('option').each(function() {
                        lhosts.push($(this).val());
                });

                if (lhosts.length < 1) {
                        $().toastmessage('showErrorToast', "No Systems To Save");
			$("launcherhost").effect("shake");
                        $("#launcherworkers").effect("shake");
                        return;
                }

		launcherconfig = new Launcher(launcherdefaultimage, lhosts);
		if (! launcherconfig.Complete()) {
			$().toastmessage('showErrorToast', "Error Saving Launcher Data!!!");
                        $("#launcherdiv").effect("shake");
                        return;
		}

		localStorage.setItem("launcherconfig", JSON.stringify(launcherconfig))

		$("#launchergrid").addClass('grid-item-complete');
                cancel_launcher();
                $().toastmessage('showSuccessToast', "Saved Launcher Information");

                if (isInstallReady()) {
                        $('#install').tooltip('hide').attr('title', 'Install Valkyrie!').tooltip('fixTitle');
                        $('#generatemanifest').tooltip('hide').attr('title', 'Generate An Installation Manifest!').tooltip('fixTitle');
                }
	}

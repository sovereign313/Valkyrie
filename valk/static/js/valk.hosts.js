	var hosts = []; 
        var hdivwidth;
        var hdivheight;

        $(function() {
		hosts = JSON.parse(localStorage.getItem("hostsconfig"));
                hdivwidth = parseInt($("#hostsdiv").css('width'), 10);
                hdivheight = parseInt($("#hostsdiv").css('height'), 10);

		if (hosts != null) {
	                $("#hostgrid").addClass('grid-item-complete');
		}

                $("#hostsdiv").hide();
                $("#hostsdiv").draggable();

		$("#hostname").keypress(function(event) {
			var keycode = (event.keyCode ? event.keyCode : event.which);
			if (keycode == '13') {
				handle_add_host();
				return;
			}
		});

		$("#listofhosts").click(function() {
			$('option:selected', this).remove();
			hosts = [];
			$('#listofhosts').find('option').each(function() {
				hosts.push($(this).val());
			});

			localStorage.removeItem("hostsconfig");
			if (hosts.length < 1) {
				return;
			}

                	localStorage.setItem("hostsconfig", JSON.stringify(hosts))
		});

                $("#hostgrid").click(function(evt) {
			offset = $("#hostgrid").offset();
			l = $("#hostgrid").offset().left - (hdivwidth / 3); 
			t = $("#hostgrid").offset().top - (hdivheight / 3) 
			$("#hostsdiv").css('left', l); 
			$("#hostsdiv").css('top', t); 

			if (hosts != null) {			
				for (i=0; i < hosts.length; i++) {
					option = document.createElement("option");
					option.text = hosts[i];
					option.value = hosts[i]; 
					document.getElementById('listofhosts').add(option);
				}
			}

			$("#hostsdiv").children().each(function() {
				$(this).css("visibility", "hidden");
			});
			$("#hostsdiv").show("scale", {}, 1000, function() {
				$("#hostsdiv").children().each(function() {
					$(this).css("visibility", "visible");
				});

				$("#hostname").focus();
			});
		});
	});

	function handle_add_host() {
		hname = document.getElementById('hostname');
		if (hname.value == "") {
			$("#hostname").effect("shake");
			$().toastmessage('showErrorToast', "The Hostname/IP box is empty");
			return;
		}

		option = document.createElement("option");
		option.text = hname.value;
		document.getElementById('listofhosts').add(option);
		hname.value="";
		hname.focus();
	}

	function cancel_hosts() {
		$("#hostsdiv").children().each(function() {
			$(this).css("visibility", "hidden");
		});

		if (hosts == null || hosts.length < 1) {
			$("#hostgrid").removeClass('grid-item-complete');
		}

		$("#hostsdiv").hide("scale", {}, 1000);
		$('#listofhosts').empty();
		$('#hostname').val('');
	}

	function save_hosts() {
		hosts = [];
		$('#listofhosts').find('option').each(function() {
			hosts.push($(this).val());
		});

		if (hosts.length < 1) {
			$().toastmessage('showErrorToast', "No Systems To Save");
			$("#listofhosts").effect("shake");
			return;
		}

                localStorage.setItem("hostsconfig", JSON.stringify(hosts))

                $("#hostgrid").addClass('grid-item-complete');
		cancel_hosts();
		$().toastmessage('showSuccessToast', "Saved Hosts List");

                if (isInstallReady()) {
                        $('#install').tooltip('hide').attr('title', 'Install Valkyrie!').tooltip('fixTitle');
                        $('#generatemanifest').tooltip('hide').attr('title', 'Generate An Installation Manifest!').tooltip('fixTitle');
                }
	}

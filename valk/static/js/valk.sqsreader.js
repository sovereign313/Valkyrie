       class SQSReader {
                constructor(sqsreadersleeptimeout, host) { 
                        this.SleepTimeout = sqsreadersleeptimeout;
			this.Host = host;
                }

                Complete() {
                        if (this.SleepTimeout == "" || this.Host == "") {
                                return false;
                        }

                        return true;
                }
        }

        var sqsreaderdivwidth;
        var sqsreaderdivheight;

        $(function() {


		sqsreaderconfig = JSON.parse(localStorage.getItem("sqsreaderconfig"));
                awsconfig = JSON.parse(localStorage.getItem("awsconfig"));

                sqsreaderdivwidth = parseInt($("#sqsreaderdiv").css('width'), 10);
                sqsreaderdivheight = parseInt($("#sqsreaderdiv").css('height'), 10);

                if (awsconfig != null) {
                        $('#awsconfiguredsqsreaderdiv').css('display', 'block');
                } else {
                        $('#awsconfiguredsqsreaderdiv').css('color', '#ff0000');
                        $('#awsconfiguredsqsreaderdiv').text("AWS IS NOT CONFIGURED!");
                        $('#awsconfiguredsqsreaderdiv').css('display', 'block');
                }

		if (sqsreaderconfig != null) {
			if (awsconfig != null) {
		                $("#sqsreadergrid").addClass('grid-item-complete');
			}
		}

                $("#sqsreaderdiv").hide();
                $("#sqsreaderdiv").draggable();

                $("#sqsreadergrid").click(function(evt) {
			offset = $("#sqsreadergrid").offset();
			l = $("#sqsreadergrid").offset().left - (sqsreaderdivwidth / 3); 
			t = $("#sqsreadergrid").offset().top - (sqsreaderdivheight / 3) 
			$("#sqsreaderdiv").css('left', l); 
			$("#sqsreaderdiv").css('top', t); 

                        awsconfig = JSON.parse(localStorage.getItem("awsconfig"));
                        if (awsconfig != null) {
                                $('#awsconfiguredsqsreaderdiv').css('display', 'block');
                                $('#awsconfiguredsqsreaderdiv').css('color', '#ffff00');
                                $('#awsconfiguredsqsreaderdiv').text("AWS Is Already Configured!");
                        } else {
                                $('#awsconfiguredsqsreaderdiv').css('color', '#ff0000');
                                $('#awsconfiguredsqsreaderdiv').text("AWS IS NOT CONFIGURED!");
                                $('#awsconfiguredsqsreaderdiv').css('display', 'block');
                        }

                        hosts = JSON.parse(localStorage.getItem("hostsconfig"));
                        $("#sqsreaderhost").empty();

                        if (hosts != null) {
                                for (i = 0; i < hosts.length; i++) {
                                        option = document.createElement("option");
                                        option.text = hosts[i];
                                        option.value = hosts[i];
                                        document.getElementById('sqsreaderhost').add(option);
                                }
                        } else {
                                option = document.createElement("option");
                                option.text = "Hosts Not Configured";
                                option.value = "none";
                                document.getElementById('sqsreaderhost').add(option);
                        }


			if (sqsreaderconfig != null) {
				$("#sqsreaderdefaultimage").value = sqsreaderconfig.DefaultImage;
                                $("#sqsreaderhost").val(sqsreaderconfig.Host);
			}

			$("#sqsreaderdiv").children().each(function() {
				$(this).css("visibility", "hidden");
			});

			$("#sqsreaderdiv").show("scale", {}, 1000, function() {
				$("#sqsreaderdiv").children().each(function() {
					$(this).css("visibility", "visible");
				});
			});
		});
	});

	function cancel_sqsreader() {
		$("#sqsreaderdiv").children().each(function() {
			$(this).css("visibility", "hidden");
		});

		$("#sqsreaderdiv").hide("scale", {}, 1000);
		$('#sqsreaderdefaultimage').val('');
	}

	function save_sqsreader() {
		sqsreadersleeptimeout = $("#sqsreadersleeptimeout").val();
                sqsreaderhost = $("#sqsreaderhost").val();

		if (sqsreadersleeptimeout == "") {
			$().toastmessage('showErrorToast', "Missing SQSReader Sleep Timeout");
                        $("#sqsreadersleeptimeout").effect("shake");
                        return;
		}

                if (sqsreaderhost == "none") {
                        $().toastmessage('showErrorToast', "Missing Host To Install On");
                        $("#sqsreaderhost").effect("shake");
                        return;
                }

		sqsreaderconfig = new SQSReader(sqsreadersleeptimeout, sqsreaderhost);
		if (! sqsreaderconfig.Complete()) {
			$().toastmessage('showErrorToast', "Error Saving SQSReader Data!!!");
                        $("#sqsreaderdiv").effect("shake");
                        return;
		}

		localStorage.setItem("sqsreaderconfig", JSON.stringify(sqsreaderconfig))
		if (awsconfig != null) {
			$("#sqsreadergrid").addClass('grid-item-complete');
		}
                cancel_sqsreader();
                $().toastmessage('showSuccessToast', "Saved SQSReader Information");

                if (isInstallReady()) {
                        $('#install').tooltip('hide').attr('title', 'Install Valkyrie!').tooltip('fixTitle');
                        $('#generatemanifest').tooltip('hide').attr('title', 'Generate An Installation Manifest!').tooltip('fixTitle');
                }
	}

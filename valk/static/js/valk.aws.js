	class AWS {
		constructor(region, sqsname, awsaccesskey, awssecretkey, awsencryptionkey) {
			this.Region = region;
			this.SQSName = sqsname;
			this.AWSAccessKey = awsaccesskey;
			this.AWSSecretKey = awssecretkey;
			this.EncryptionKey = awsencryptionkey; 
		}

		Complete() {
			if (this.Region == "" || this.SQSName == "" || this.AWSAccessKey == "" || this.AWSSecretKey == "" || this.EncryptionKey == "") {
				return false;
			}

			return true;
		}
	}

        $(function() {

	        awsconfig = JSON.parse(localStorage.getItem("awsconfig"));
		awsdivwidth = parseInt($("#awsdiv").css('width'), 10);
                awsdivheight = parseInt($("#awsdiv").css('height'), 10);

		if (awsconfig != null) {
			$("#awsgrid").addClass('grid-item-complete');
		}

		$("#awsdiv").hide();
		$("#awsdiv").draggable();
	
		$("#awsgrid").click(function(evt) {
	                offset = $("#awsgrid").offset();
	                l = $("#awsgrid").offset().left - (awsdivwidth / 3);
        	        t = $("#awsgrid").offset().top - (awsdivheight / 3) 
                	$("#awsdiv").css('left', l);
	                $("#awsdiv").css('top', t);

			if (awsconfig != null) {
				$("#awsregion").val(awsconfig.Region);
				$("#awssqsname").val(awsconfig.SQSName);
				$("#awsaccesskey").val(awsconfig.AWSAccessKey);
				$("#awssecretkey").val(awsconfig.AWSSecretKey);
				$("#awsencryptionkey").val(awsconfig.EncryptionKey);
			}
	
			$("#awsdiv").children().each(function() {
				$(this).css("visibility", "hidden");
			});

			$("#awsdiv").show("scale", {}, 1000, function() {
				$("#awsdiv").children().each(function() {
					$(this).css("visibility", "visible");
				});

				$("#awsregion").focus();
			});
		});

		$("#awsregion").keypress(function(event) {
			var keycode = (event.keyCode ? event.keyCode : event.which);
			if (keycode == '13') {
				flag = true;
				$("#awsdiv").find('input[type="text"]').each(function() {
					if ($(this).val() == "") {
						flag = false;
					}
				});
				
				if (! flag) {
					$("#awssqsname").focus();
				} else {
					save_aws();
				}

				return;
			}
		});

		$("#awssqsname").keypress(function(event) {
			var keycode = (event.keyCode ? event.keyCode : event.which);
			if (keycode == '13') {
				flag = true;
				$("#awsdiv").find('input[type="text"]').each(function() {
					if ($(this).val() == "") {
						flag = false;
					}
				});
				
				if (! flag) {
					$("#awsaccesskey").focus();
				} else {
					save_aws();
				}
				return;
			}
		});

		$("#awsaccesskey").keypress(function(event) {
			var keycode = (event.keyCode ? event.keyCode : event.which);
			if (keycode == '13') {
				flag = true;
				$("#awsdiv").find('input[type="text"]').each(function() {
					if ($(this).val() == "") {
						flag = false;
					}
				});
				
				if (! flag) {
					$("#awssecretkey").focus();
				} else {
					save_aws();
				}

				return;
			}
		});

		$("#awssecretkey").keypress(function(event) {
			var keycode = (event.keyCode ? event.keyCode : event.which);
			if (keycode == '13') {
				flag = true;
				$("#awsdiv").find('input[type="text"]').each(function() {
					if ($(this).val() == "") {
						flag = false;
					}
				});
				
				if (! flag) {
					$("#awsencryptionkey").focus();
				} else {
					save_aws();
				}

				return;
			}
		});

		$("#awsencryptionkey").keypress(function(event) {
			var keycode = (event.keyCode ? event.keyCode : event.which);
			if (keycode == '13') {
				save_aws();
				return;
			}
		});

	});

        function cancel_aws() {
                $("#awsdiv").children().each(function() {
                        $(this).css("visibility", "hidden");
                });

		$("#awsdiv").find('input[type="text"]').val('');
                $("#awsdiv").hide("scale", {}, 1000);
        }

        function save_aws() {
		region = $("#awsregion").val();
		sqsname = $("#awssqsname").val();
		awsaccesskey = $("#awsaccesskey").val();
		awssecretkey = $("#awssecretkey").val();
		awsencryptionkey = $("#awsencryptionkey").val();

		if (region == "") {
			$().toastmessage('showErrorToast', "Missing AWS Region");
                        $("#awsregion").effect("shake");
                        return;
		}

		if (sqsname == "") {
			$().toastmessage('showErrorToast', "Missing AWS SQS Name");
                        $("#awssqsname").effect("shake");
                        return;
		}

		if (awsaccesskey == "") {
			$().toastmessage('showErrorToast', "Missing AWS Access Key");
                        $("#awsaccesskey").effect("shake");
                        return;
		}

		if (awssecretkey == "") {
			$().toastmessage('showErrorToast', "Missing AWS Secret Key");
                        $("#awssecretkey").effect("shake");
                        return;
		}

		if (awsencryptionkey == "") {
			$().toastmessage('showErrorToast', "Missing AWS Encryption Key");
                        $("#awsencryptionkey").effect("shake");
                        return;
		}

		if (awsencryptionkey == "LKHlhb899Y09olUi") {
			$().toastmessage('showErrorToast', "That Key Is Public, And Will Get You Hacked!  Please Select Another");
                        $("#awsencryptionkey").effect("shake");
                        return;
		}

		awsconfig = new AWS(region, sqsname, awsaccesskey, awssecretkey, awsencryptionkey);
		if (! awsconfig.Complete()) {
			$().toastmessage('showErrorToast', "Error Saving AWS Data!!!");
                        $("#awsdiv").effect("shake");
                        return;
		}

		localStorage.setItem("awsconfig", JSON.stringify(awsconfig))
		$("#awsgrid").addClass('grid-item-complete');
	
		gnuff = localStorage.getItem("alerterconfig");
		if (gnuff == null) {
			$("#alertergrid").addClass('grid-item-goodenough');
		}


                cancel_aws();
                $().toastmessage('showSuccessToast', "Saved AWS Information");

                if (isInstallReady()) {
                        $('#install').tooltip('hide').attr('title', 'Install Valkyrie!').tooltip('fixTitle');
                        $('#generatemanifest').tooltip('hide').attr('title', 'Generate An Installation Manifest!').tooltip('fixTitle');
                }
        }


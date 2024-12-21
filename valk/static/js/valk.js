	var licensekey;
	var business;
	var refreshId;

        $(function() {
		manifestdivwidth = parseInt($("#manifestdiv").css('width'), 10);
		manifestdivheight = parseInt($("#manifestdiv").css('height'), 10);
		updatelogdivwidth = parseInt($("#updatelogdiv").css('width'), 10);
		updatelogdivheight = parseInt($("#updatelogdiv").css('height'), 10);

		$('[data-toggle="tooltip"]').tooltip();
		$("#mainlist").hide();
		$("#showoptions").hide();
		$("#install").hide();
		$("#generatemanifest").hide();
		$("#clearsetup").hide();
		$("#licensediv").hide();
		$("#manifestdiv").hide();
		$("#updatelogdiv").hide();

		$("#updatelogdiv").draggable();

		licensekey = localStorage.getItem("licensekey");
		business = localStorage.getItem("business");
		if (licensekey == null || business == null) {
			$('#install').tooltip('hide').attr('title', 'A Valid License Is Required To Proceed!').tooltip('fixTitle');
			$('#generatemanifest').tooltip('hide').attr('title', 'A Valid License Is Required To Proceed!').tooltip('fixTitle');
		}

		window.setTimeout(function() {
			$("#mainimg").slideUp('slow', function() {
				$("#mainlist").fadeIn('slow');
				$("#install").fadeIn('slow');
				$("#generatemanifest").fadeIn('slow');
				$("#clearsetup").fadeIn('slow');
				$("#showoptions").slideDown('fast');
				$("#showoptions").val("local");
				if (licensekey == null || business == null) {
					$("#licensediv").slideDown('slow');
				}

				if (isInstallReady()) {
					$('#install').tooltip('hide').attr('title', 'Install Valkyrie!').tooltip('fixTitle');
					$('#generatemanifest').tooltip('hide').attr('title', 'Generate An Installation Manifest!').tooltip('fixTitle');
				}

				show_specific();
			});
		}, 1500);

		$("#logodiv").click(function() {
			if ($("#licensediv").is(":visible")) {
				$("#licensediv").slideUp('fast');
			} else {
				$("#license").val(licensekey);
				$("#business").val(business);
				$("#licensediv").slideDown('fast');
			}
		});

                $("#license").keypress(function(event) {
                        var keycode = (event.keyCode ? event.keyCode : event.which);
                        if (keycode == '13') {
				licensekey = $("#license").val();
				business = $("#business").val();
				if (licensekey == "") {
					$("#license").effect("shake");
					$().toastmessage('showErrorToast', "Missing License Key");
					return;
				}

				if (business == "") {
					$("#business").effect("shake");
					$().toastmessage('showErrorToast', "Missing Business Name");
					return;
				}

				var data = {
					licensekey:licensekey,
					business:business
				}

				$.post(
					'/verifykey',
					data,
					function(responseText) {
						resp = JSON.parse(responseText);
						if (resp.Code == "200") {
							localStorage.setItem("licensekey", licensekey);
							localStorage.setItem("business", business);
							$("#licensediv").slideUp('fast');
							$().toastmessage('showSuccessToast', "Activiated Valkyrie");

							if (isInstallReady()) {
								$('#install').tooltip('hide').attr('title', 'Install Valkyrie!').tooltip('fixTitle');
								$('#generatemanifest').tooltip('hide').attr('title', 'Generate An Installation Manifest!').tooltip('fixTitle');
							} else {
								$('#install').tooltip('hide').attr('title', 'Valkyrie Needs To Be Configured!').tooltip('fixTitle');
								$('#generatemanifest').tooltip('hide').attr('title', 'Valkyrie Needs To Be Configured!').tooltip('fixTitle');
							}

							return;
						} else if (resp.Code == "502") {
							$().toastmessage('showErrorToast', "Tampering Has Been Detected.  Please Contact Support");
							return;
						} else {
							$("#licensediv").effect("shake");
							$().toastmessage('showErrorToast', "License Key Is Not Valid!");
							return;
						}
					}
				).fail(function(response) {
					$().toastmessage('showErrorToast', 'Failed To Contact Installer Backend. Is It Still Running?: ' + response);
					return;
				});
                        }
                });

                $("#business").keypress(function(event) {
                        var keycode = (event.keyCode ? event.keyCode : event.which);
                        if (keycode == '13') {
				licensekey = $("#license").val();
				business = $("#business").val();
				if (licensekey == "") {
					$("#license").effect("shake");
					$().toastmessage('showErrorToast', "Missing License Key");
					return;
				}

				if (business == "") {
					$("#business").effect("shake");
					$().toastmessage('showErrorToast', "Missing Business Name");
					return;
				}

				var data = {
					licensekey:licensekey,
					business:business
				}

				$.post(
					'/verifykey',
					data,
					function(responseText) {
						resp = JSON.parse(responseText);
						if (resp.Code == "200") {
							localStorage.setItem("licensekey", licensekey);
							localStorage.setItem("business", business);
							$("#licensediv").slideUp('fast');
							$().toastmessage('showSuccessToast', "Activiated Valkyrie");

							if (isInstallReady()) {
								$('#install').tooltip('hide').attr('title', 'Install Valkyrie!').tooltip('fixTitle');
								$('#generatemanifest').tooltip('hide').attr('title', 'Generate An Installation Manifest!').tooltip('fixTitle');
							} else {
								$('#install').tooltip('hide').attr('title', 'Valkyrie Needs To Be Configured!').tooltip('fixTitle');
								$('#generatemanifest').tooltip('hide').attr('title', 'Valkyrie Needs To Be Configured!').tooltip('fixTitle');
							}

							return;
						} else if (resp.Code == "502") {
							$().toastmessage('showErrorToast', "Tampering Has Been Detected.  Please Contact Support");
							return;
						} else {
							$("#licensediv").effect("shake");
							$().toastmessage('showErrorToast', "License Key Is Not Valid!");
							return;
						}
					}
				).fail(function(response) {
					$().toastmessage('showErrorToast', 'Failed To Contact Installer Backend. Is It Still Running?: ' + response);
					return;
				});
                        }
                });

		$("#install").mouseover(function() {
			if (isInstallReady()) {
				$(this).addClass("active");
				return;
			}
		});

		$("#generatemanifest").mouseover(function() {
			if (isInstallReady()) {
				$(this).addClass("active");
				return;
			}
		});

		$("#information").click(function() {
			window.open("http://wiki.valkyriesoftware.io");
		});

		$("#install").click(function() {
			if (! isInstallReady()) {
				$(this).effect("shake");
				$().toastmessage('showErrorToast', "Something isn't configured");
				return;
			}

			license = localStorage.getItem("licensekey");
			business = localStorage.getItem("business");
			hc = localStorage.getItem("hostsconfig");
			fc = localStorage.getItem("foremanconfig");
			lc = localStorage.getItem("launcherconfig");
			wc = localStorage.getItem("workerconfig");
			lgc = localStorage.getItem("loggerconfig");
			ac = localStorage.getItem("alerterconfig");
			awc = localStorage.getItem("awsconfig");
			dc = localStorage.getItem("dispatcherconfig");
			sc = localStorage.getItem("sqsreaderconfig");
			mc = localStorage.getItem("mailreaderconfig");

			if (hc == null) {
				hc = "none";
			}

			if (fc == null) {
				fc = "none";
			}

			if (lc == null) {
				lc = "none";
			}

			if (wc == null) {
				wc = "none";
			}

			if (lgc == null) {
				lgc = "none";
			}

			if (ac == null) {
				ac = "none";
			}

			if (awc == null) {
				awc = "none";
			}

			if (dc == null) {
				dc = "none";
			}

			if (sc == null) {
				sc = "none";
			}

			if (mc == null) {
				mc = "none";
			}

			var data = {
				licensekey:licensekey,
				business:business,
				hostconfig:hc,
				foremanconfig:fc,
				launcherconfig:lc,
				workerconfig:wc,
				loggerconfig:lgc,
				alerterconfig:ac,
				awsconfig:awc,
				dispatcherconfig:dc,
				sqsreaderconfig:sc,
				mailreaderconfig:mc
			}

			$("#install").slideUp('fast', function() {
				$("#spinner").slideDown('fast');
			});

			offset = $("#install").offset();
			l = $("#install").offset().left - (updatelogdivwidth / 3);
			t = $("#install").offset().top - (updatelogdivheight / 3)
			$("#updatelogdiv").css('left', l);
			$("#updatelogdiv").css('top', t - 250);

			$("#updatelogdiv").children().each(function() {
				$(this).css("visibility", "hidden");
			});

			$("#updatelogdiv").show("scale", {}, 1000, function() {
				$("#updatelogdiv").children().each(function() {
					$(this).css("visibility", "visible");
				});
			});

			window.setTimeout(function() {
				$().toastmessage('showNoticeToast', "Seriously, Get A Coffee Or Something.  This takes a while.");
				window.setTimeout(function() {
					$().toastmessage('showNoticeToast', "Maybe Check Out this web comic: https://xkcd.com/1319/");
					window.setTimeout(function() {
						$().toastmessage('showNoticeToast', "Ok, Ok. A Bit, A Byte, And Nibble Walk Into A Bar...");
						window.setTimeout(function() {
							$().toastmessage('showNoticeToast', "The Bartender says \"What's The Word?\"");
						}, 10000);
					}, 60000);
				}, 60000);
			}, 45000);

			refreshId = setInterval(function() {
				$.get('/status', function(data, status) {
					$("#logdata").val(data);
					$("#logdata").scrollTop($("#logdata").scrollHeight);
				});
			}, 2000);

			$.post(
			'/install',
			data,
			function(responseText) {
				resp = JSON.parse(responseText);
				$("#spinner").slideUp('fast', function() {
					$("#install").slideDown('fast');
				});

				if (resp.Code == "200") {
					$().toastmessage('showSuccessToast', resp.Message);
					$("#spinner").stop(true);
					$("#install").stop(true);
					$("#spinner").slideUp('fast', function() {
						$("#install").slideDown('fast');
					});
					return;
				} else if (resp.Code == "502") {
					$().toastmessage('showErrorToast', "Tampering Has Been Detected.  Please Contact Support");
					$("#spinner").stop(true);
					$("#install").stop(true);
					$("#spinner").slideUp('fast', function() {
						$("#install").slideDown('fast');
					});
					return;
				} else if (resp.Code == "501") {
					$("#licensediv").effect("shake");
					$().toastmessage('showErrorToast', "License Key Is Not Valid!");
					return
				} else {
					$().toastmessage('showErrorToast', "Error Code: " + resp.Code + " Message: " + resp.Message);
					return;
				}
			}
			).fail(function(response) {
				$("#spinner").stop(true);
				$("#install").stop(true);
				$("#spinner").slideUp('fast', function() {
					$("#install").slideDown('fast');
				});
				$("#updatelogdiv").stop(true);
				$("#updatelogdiv").hide();
				$().toastmessage('showErrorToast', 'Failed To Contact Installer Backend. Is It Still Running?');
				return;
			});

		});

		$("#generatemanifest").click(function() {
			if (! isInstallReady()) {
				$(this).effect("shake");
				$().toastmessage('showErrorToast', "Something isn't configured");
				return;
			}

			license = localStorage.getItem("licensekey");
			business = localStorage.getItem("business");
			hc = localStorage.getItem("hostsconfig");
			fc = localStorage.getItem("foremanconfig");
			lc = localStorage.getItem("launcherconfig");
			wc = localStorage.getItem("workerconfig");
			lgc = localStorage.getItem("loggerconfig");
			ac = localStorage.getItem("alerterconfig");
			awc = localStorage.getItem("awsconfig");
			dc = localStorage.getItem("dispatcherconfig");
			sc = localStorage.getItem("sqsreaderconfig");
			mc = localStorage.getItem("mailreaderconfig");

			if (hc == null) {
				hc = "none";
			}

			if (fc == null) {
				fc = "none";
			}

			if (lc == null) {
				lc = "none";
			}

			if (wc == null) {
				wc = "none";
			}

			if (lgc == null) {
				lgc = "none";
			}

			if (ac == null) {
				ac = "none";
			}

			if (awc == null) {
				awc = "none";
			}

			if (dc == null) {
				dc = "none";
			}

			if (sc == null) {
				sc = "none";
			}

			if (mc == null) {
				mc = "none";
			}

			var data = {
				licensekey:licensekey,
				business:business,
				hostconfig:hc,
				foremanconfig:fc,
				launcherconfig:lc,
				workerconfig:wc,
				loggerconfig:lgc,
				alerterconfig:ac,
				awsconfig:awc,
				dispatcherconfig:dc,
				sqsreaderconfig:sc,
				mailreaderconfig:mc
			}

			$.post(
			'/generatemanifest',
			data,
			function(responseText) {
				resp = JSON.parse(responseText);
				if (resp.Code == "200") {
					$("#manifest").val(resp.Message);
                		        offset = $("#generatemanifest").offset();
		                        l = $("#generatemanifest").offset().left - (manifestdivwidth / 3);
		                        t = $("#generatemanifest").offset().top - (manifestdivheight / 3)
                		        $("#manifestdiv").css('left', l);
		                        $("#manifestdiv").css('top', t - 150);
		
		                        $("#manifestdiv").children().each(function() {
                		                $(this).css("visibility", "hidden");
		                        });

                		        $("#manifestdiv").show("scale", {}, 1000, function() {
                                		$("#manifestdiv").children().each(function() {
		                                        $(this).css("visibility", "visible");
                		                });
		                        });
					return;
				} else if (resp.Code == "502") {
					$().toastmessage('showErrorToast', "Tampering Has Been Detected.  Please Contact Support");
					return;
				} else {
					$("#licensediv").effect("shake");
					$().toastmessage('showErrorToast', "License Key Is Not Valid!");
					return;
				}
			}
			).fail(function(response) {
				$().toastmessage('showErrorToast', 'Failed To Contact Installer Backend. Is It Still Running?: ' + response);
				return;
			});
		
		
		});

		$("#clearsetup").click(function() {
			$.confirm({
				title: 'Woah!',
				content: 'Really Delete All Your Setup?',
				buttons: {
					confirm: function() {
						localStorage.removeItem("hostsconfig");
						localStorage.removeItem("foremanconfig");
						localStorage.removeItem("launcherconfig");
						localStorage.removeItem("workerconfig");
						localStorage.removeItem("loggerconfig");
						localStorage.removeItem("alerterconfig");
						localStorage.removeItem("awsconfig");
						localStorage.removeItem("dispatcherconfig");
						localStorage.removeItem("sqsreaderconfig");
						localStorage.removeItem("mailreaderconfig");
						localStorage.removeItem("licensekey");
						localStorage.removeItem("business");
						$().toastmessage('showSuccessToast', "Cleared The Setup!");
						window.location.reload(true);
					},
					cancel: function() {
						$().toastmessage('showNoticeToast', "Whew! Not Deleting Anything");
					}
				}
			});	
		});
	});

	function save_manifest() {
		$("<a />", {
			download: $.now() + ".json",
			href: URL.createObjectURL(
				new Blob([$('#manifest').val()], {
					type: "text/plain"
				}))
			
		}).appendTo("body")[0].click();

		$(window).one("focus", function() {
			$("a").last().remove();
		});
	}

	function cancel_manifest() {
		$("#manifestdiv").children().each(function() {
			$(this).css("visibility", "hidden");
                });

                $("#manifestdiv").hide("scale", {}, 1000);
                $('#manifest').val('');
	}

	function cancel_logdiv() {
		clearInterval(refreshId);
		$("#updatelogdiv").children().each(function() {
			$(this).css("visibility", "hidden");
                });

                $("#updatelogdiv").hide("scale", {}, 1000);
                $('#logdata').val('');
	}

	function isInstallReady() {
		hostc = JSON.parse(localStorage.getItem("hostsconfig"));
		foremanc = JSON.parse(localStorage.getItem("foremanconfig"));
		launcherc = JSON.parse(localStorage.getItem("launcherconfig"));
		workerc = JSON.parse(localStorage.getItem("workerconfig"));
		loggerc = JSON.parse(localStorage.getItem("loggerconfig"));
		alerterc = JSON.parse(localStorage.getItem("alerterconfig"));
		awsc = JSON.parse(localStorage.getItem("awsconfig"));
		dispatcherc = JSON.parse(localStorage.getItem("dispatcherconfig"));
		sqsreaderc = JSON.parse(localStorage.getItem("sqsreaderconfig"));
		mailreaderc = JSON.parse(localStorage.getItem("mailreaderconfig"));

		if (licensekey == null) {
			return false;
		}

		if (business == null) {
			return false
		}

		if (hostc != null && foremanc != null && launcherc != null && workerc != null) {
			return true;
		}

		if (hostc != null && launcherc != null && workerc != null && awsc != null && dispatcherc != null && sqsreader != null) {
			return true;
		}

		return false;
	}

	function show_specific() {
		showoption = $('#showoptions').val();
		if (showoption == "all") {
			$('#hostgrid').slideDown('slow');
			$('#awsgrid').slideDown('slow');
			$('#alertergrid').slideDown('slow');
			$('#loggergrid').slideDown('slow');
			$('#mailreadergrid').slideDown('slow');
			$('#foremangrid').slideDown('slow');
			$('#launchergrid').slideDown('slow');
			$('#workergrid').slideDown('slow');
			$('#dispatchergrid').slideDown('slow');
			$('#sqsreadergrid').slideDown('slow');
			$('#informationgrid').slideDown('slow');
			$('#supportgrid').slideDown('slow');
		} else if (showoption == "local") {
			$('#awsgrid').slideUp('slow');
			$('#dispatchergrid').slideUp('slow');
			$('#sqsreadergrid').slideUp('slow');
			$('#mailreadergrid').slideUp('slow');
			$('#hostgrid').slideDown('slow');
			$('#foremangrid').slideDown('slow');
			$('#launchergrid').slideDown('slow');
			$('#workergrid').slideDown('slow');
			$('#alertergrid').slideDown('slow');
			$('#loggergrid').slideDown('slow');
			$('#informationgrid').slideDown('slow');
			$('#supportgrid').slideDown('slow');
		} else if (showoption == "aws") {
			$('#awsgrid').slideDown('slow');
			$('#dispatchergrid').slideDown('slow');
			$('#sqsreadergrid').slideDown('slow');
			$('#foremangrid').slideUp('slow');
			$('#hostgrid').slideDown('slow');
			$('#launchergrid').slideDown('slow');
			$('#workergrid').slideDown('slow');
			$('#alertergrid').slideDown('slow');
			$('#loggergrid').slideDown('slow');
			$('#mailreadergrid').slideUp('slow');
			$('#informationgrid').slideDown('slow');
			$('#supportgrid').slideDown('slow');
		} else if (showoption == "localreq") {
			$('#hostgrid').slideDown('slow');
			$('#foremangrid').slideDown('slow');
			$('#launchergrid').slideDown('slow');
			$('#workergrid').slideDown('slow');
			$('#awsgrid').slideUp('slow');
			$('#alertergrid').slideUp('slow');
			$('#loggergrid').slideUp('slow');
			$('#mailreadergrid').slideUp('slow');
			$('#dispatchergrid').slideUp('slow');
			$('#sqsreadergrid').slideUp('slow');
			$('#informationgrid').slideUp('slow');
			$('#supportgrid').slideUp('slow');
		} else if (showoption == "awsreq") {
			$('#hostgrid').slideDown('slow');
			$('#foremangrid').slideUp('slow');
			$('#launchergrid').slideDown('slow');
			$('#workergrid').slideDown('slow');
			$('#awsgrid').slideDown('slow');
			$('#alertergrid').slideUp('slow');
			$('#loggergrid').slideUp('slow');
			$('#mailreadergrid').slideUp('slow');
			$('#dispatchergrid').slideDown('slow');
			$('#sqsreadergrid').slideDown('slow');
			$('#informationgrid').slideUp('slow');
			$('#supportgrid').slideUp('slow');
		}
	}



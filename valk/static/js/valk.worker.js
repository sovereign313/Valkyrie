       class Worker {
                constructor(externalpath, repourl, sshuser, privatekey, publickey) { 
			this.ExternalPath = externalpath;
			this.GitRepoUrl = repourl; 
			this.SSHUser = sshuser;
			this.SSHPrivateKey = privatekey;
			this.SSHPublicKey = publickey;
                }

                Complete() {
                        if (this.ExternalPath == "" || this.GitRepoUrl == "" || this.SSHUser == "" || this.SSHPrivateKey == "" || this.SSHPublicKey == "") {
                                return false;
                        }

                        return true;
                }
        }

        var workerdivwidth;
        var workerdivheight;

        $(function() {

		workerconfig = JSON.parse(localStorage.getItem("workerconfig"));
                workerdivwidth = parseInt($("#workerdiv").css('width'), 10);
                workerdivheight = parseInt($("#workerdiv").css('height'), 10);

		if (workerconfig != null) {
	                $("#workergrid").addClass('grid-item-complete');
		}

                $("#workerdiv").hide();
                $("#workerdiv").draggable();

                $("#workergrid").click(function(evt) {
			offset = $("#workergrid").offset();
			l = $("#workergrid").offset().left - (workerdivwidth / 3); 
			t = $("#workergrid").offset().top - (workerdivheight / 3) 
			$("#workerdiv").css('left', l); 
			$("#workerdiv").css('top', t);

			if (workerconfig != null) {
				$("#workerexternalpath").val(workerconfig.ExternalPath);
				$("#workergitrepourl").val(workerconfig.GitRepoUrl);
				$("#workersshuser").val(workerconfig.SSHUser);
				$("#workerprivatekey").val(workerconfig.SSHPrivateKey);
				$("#workerpublickey").val(workerconfig.SSHPublicKey);
			}

			$("#workerdiv").children().each(function() {
				$(this).css("visibility", "hidden");
			});

			$("#workerdiv").show("scale", {}, 1000, function() {
				$("#workerdiv").children().each(function() {
					$(this).css("visibility", "visible");
				});
			});
		});
	});

	function cancel_worker() {
		$("#workerdiv").children().each(function() {
			$(this).css("visibility", "hidden");
		});

		$("#workerdiv").hide("scale", {}, 1000);
		$('#workerdefaultimage').val('');
	}

	function save_worker() {
		externalpath = $("#workerexternalpath").val();
		repourl = $("#workergitrepourl").val();
		sshuser = $("#workersshuser").val();
		privkey = $("#workerprivatekey").val();
		pubkey = $("#workerpublickey").val();

		if (externalpath == "") {
			$().toastmessage('showErrorToast', "Missing External Path");
                        $("#workerexternalpath").effect("shake");
                        return;
		}

		if (repourl == "") {
			$().toastmessage('showErrorToast', "Missing Repo URL");
                        $("#workerrepourl").effect("shake");
                        return;
		}

		if (sshuser == "") {
			$().toastmessage('showErrorToast', "Missing SSH User");
                        $("#workersshuser").effect("shake");
                        return;
		}

		if (privkey == "") {
			$().toastmessage('showErrorToast', "Missing SSH Private Key");
			$("#workerprivatekey").effect("shake");
			return;
		}

		if (pubkey == "") {
			$().toastmessage('showErrorToast', "Missing SSH Public Key");
			$("#workerpublickey").effect("shake");
			return;
		}

		if (sshuser != "root") {
			$().toastmessage('showNoticeToast', "It is recommended that you run as root");
		}

		workerconfig = new Worker(externalpath, repourl, sshuser, privkey, pubkey);
		if (! workerconfig.Complete()) {
			$().toastmessage('showErrorToast', "Error Saving Worker Data!!!");
                        $("#workerdiv").effect("shake");
                        return;
		}

		localStorage.setItem("workerconfig", JSON.stringify(workerconfig))
		$("#workergrid").addClass('grid-item-complete');
                cancel_worker();
                $().toastmessage('showSuccessToast', "Saved Worker Information");

                if (isInstallReady()) {
                        $('#install').tooltip('hide').attr('title', 'Install Valkyrie!').tooltip('fixTitle');
                        $('#generatemanifest').tooltip('hide').attr('title', 'Generate An Installation Manifest!').tooltip('fixTitle');
                }
	}

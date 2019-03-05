/*!
 * Image (upload) dialog plugin for Editor.md
 *
 * @file        image-dialog.js
 * @author      pandao
 * @version     1.3.4
 * @updateTime  2015-06-09
 * {@link       https://github.com/pandao/editor.md}
 * @license     MIT
 */

(function () {

    var factory = function (exports) {

        var pluginName = "image-dialog";

        exports.fn.imageDialog = function () {

            var _this = this;
            var cm = this.cm;
            var lang = this.lang;
            var editor = this.editor;
            var settings = this.settings;
            var cursor = cm.getCursor();
            var selection = cm.getSelection();
            var imageLang = lang.dialog.image;
            var classPrefix = this.classPrefix;
            var iframeName = classPrefix + "image-iframe";
            var dialogName = classPrefix + pluginName, dialog;

            cm.focus();

            var loading = function (show) {
                var _loading = dialog.find("." + classPrefix + "dialog-mask");
                _loading[(show) ? "show" : "hide"]();
            };

            if (editor.find("." + dialogName).length < 1) {
                var guid = (new Date).getTime();
                var action = settings.imageUploadURL



                var dialogContent = ((settings.imageUpload) ? "<form id=\"update\" action=\"" + action + "\" method=\"post\" enctype=\"multipart/form-models\" class=\"" + classPrefix + "form\">" : "<div class=\"" + classPrefix + "form\">") +
                    ((settings.imageUpload) ? "<iframe name=\"" + iframeName + "\" id=\"" + iframeName + "\" guid=\"" + guid + "\"></iframe>" : "") +
                    "<label>" + imageLang.url + "</label>" +
                    "<input id=\"imgurlt\" type=\"text\" models-url />" + (function () {
                        return (settings.imageUpload) ? "<div class=\"" + classPrefix + "file-input\">" +
                            "<input type=\"file\" name=\"smfile\" accept=\"image/*\" />" +
                            "<input type=\"button\" class=\"editormd-btn\" value=\"" + imageLang.uploadButton + "\" />" +
                            "</div>" : "";
                    })() +
                    "<br/>" +
                    "<label>" + imageLang.alt + "</label>" +
                    "<input type=\"text\" value=\"" + selection + "\" models-alt />" +
                    "<br/>" +
                    "<label>" + imageLang.link + "</label>" +
                    "<input id=\"imgurl\" type=\"text\" value=\"http://\" models-link />" +
                    "<br/>" +
                    ((settings.imageUpload) ? "</form>" : "</div>");

                //var imageFooterHTML = "<button class=\"" + classPrefix + "btn " + classPrefix + "image-manager-btn\" style=\"float:left;\">" + imageLang.managerButton + "</button>";

                dialog = this.createDialog({
                    title: imageLang.title,
                    width: (settings.imageUpload) ? 465 : 380,
                    height: 254,
                    name: dialogName,
                    content: dialogContent,
                    mask: settings.dialogShowMask,
                    drag: settings.dialogDraggable,
                    lockScreen: settings.dialogLockScreen,
                    maskStyle: {
                        opacity: settings.dialogMaskOpacity,
                        backgroundColor: settings.dialogMaskBgColor
                    },
                    buttons: {
                        enter: [lang.buttons.enter, function () {
                            var url = this.find("[models-url]").val();
                            var alt = this.find("[models-alt]").val();
                            var link = this.find("[models-link]").val();

                            if (url === "") {
                                alert(imageLang.imageURLEmpty);
                                return false;
                            }

                            var altAttr = (alt !== "") ? " \"" + alt + "\"" : "";

                            if (link === "" || link === "http://") {
                                cm.replaceSelection("![" + alt + "](" + url + altAttr + ")");
                            }
                            else {
                                cm.replaceSelection("[![" + alt + "](" + url + altAttr + ")](" + link + altAttr + ")");
                            }

                            if (alt === "") {
                                cm.setCursor(cursor.line, cursor.ch + 2);
                            }

                            this.hide().lockScreen(false).hideMask();

                            return false;
                        }],

                        cancel: [lang.buttons.cancel, function () {
                            this.hide().lockScreen(false).hideMask();

                            return false;
                        }]
                    }
                });

                dialog.attr("id", classPrefix + "image-dialog-" + guid);

                if (!settings.imageUpload) {
                    return;
                }

                var fileInput = dialog.find("[name=\"smfile\"]");

                fileInput.bind("change", function () {
                    var fileName = fileInput.val();
                    var isImage = new RegExp("(\\.(" + settings.imageFormats.join("|") + "))$"); // /(\.(webp|jpg|jpeg|gif|bmp|png))$/
                    
                    if (fileName === "") {
                        alert(imageLang.uploadFileEmpty);

                        return false;
                    }

                    if (!isImage.test(fileName)) {
                        alert(imageLang.formatNotAllowed + settings.imageFormats.join(", "));

                        return false;
                    }

                    loading(true);
                    update()

                });
            }

            dialog = editor.find("." + dialogName);
            dialog.find("[type=\"text\"]").val("");
            dialog.find("[type=\"file\"]").val("");
            dialog.find("[models-link]").val("http://");

            this.dialogShowMask(dialog);
            this.dialogLockScreen();
            dialog.show();

        };

    };

    // CommonJS/Node.js
    if (typeof require === "function" && typeof exports === "object" && typeof module === "object") {
        module.exports = factory;
    }
    else if (typeof define === "function")  // AMD/CMD/Sea.js
    {
        if (define.amd) { // for Require.js

            define(["editormd"], function (editormd) {
                factory(editormd);
            });

        } else { // for Sea.js
            define(function (require) {
                var editormd = require("./../../editormd");
                factory(editormd);
            });
        }
    }
    else {
        factory(window.editormd);
    }

})();


function update() {
    function showResponse(responseText, statusText) {
        console.log(responseText)
        $("#imgurlt").val(responseText.data.url)
        $(".editormd-dialog-mask").hide()
    } 
    var options = {
        success:showResponse
    }
    $("#update").ajaxSubmit(options);

}


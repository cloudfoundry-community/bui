ace.define("ace/theme/panda",["require","exports","module","ace/lib/dom"], function(require, exports, module) {

exports.isDark = true;
exports.cssClass = "ace-panda";
exports.cssText = ".ace-panda .ace_gutter {\
background: #292A2B;\
color: rgb(136,136,137)\
}\
.ace-panda .ace_print-margin {\
width: 1px;\
background: #e8e8e8\
}\
.ace-panda {\
background-color: #292A2B;\
color: #E6E6E6\
}\
.ace-panda .ace_cursor {\
color: #F8F8F0\
}\
.ace-panda .ace_marker-layer .ace_selection {\
background: rgba(60, 64, 69, 0.53)\
}\
.ace-panda.ace_multiselect .ace_selection.ace_start {\
box-shadow: 0 0 3px 0px #292A2B;\
border-radius: 2px\
}\
.ace-panda .ace_marker-layer .ace_step {\
background: rgb(198, 219, 174)\
}\
.ace-panda .ace_marker-layer .ace_bracket {\
margin: -1px 0 0 -1px;\
border: 1px solid #34383D\
}\
.ace-panda .ace_marker-layer .ace_active-line {\
background: #404954\
}\
.ace-panda .ace_gutter-active-line {\
background-color: #404954\
}\
.ace-panda .ace_marker-layer .ace_selected-word {\
border: 1px solid rgba(60, 64, 69, 0.53)\
}\
.ace-panda .ace_fold {\
background-color: #65BDFF;\
border-color: #E6E6E6\
}\
.ace-panda .ace_keyword {\
color: #FF75B5\
}\
.ace-panda .ace_constant.ace_language,\
.ace-panda .ace_constant.ace_numeric,\
.ace-panda .ace_entity.ace_other.ace_attribute-name,\
.ace-panda .ace_storage,\
.ace-panda .ace_storage.ace_type,\
.ace-panda .ace_support.ace_constant {\
color: #FFB86C\
}\
.ace-panda .ace_constant.ace_character,\
.ace-panda .ace_constant.ace_other {\
color: #6DB1FF\
}\
.ace-panda .ace_entity.ace_name.ace_function,\
.ace-panda .ace_support.ace_function {\
color: #65BDFF\
}\
.ace-panda .ace_support.ace_class,\
.ace-panda .ace_support.ace_type {\
color: #FFC990\
}\
.ace-panda .ace_invalid {\
color: #EBEBEB\
}\
.ace-panda .ace_invalid.ace_deprecated {\
text-decoration: underline;\
font-style: italic;\
color: #7A6E71\
}\
.ace-panda .ace_string {\
color: #19F9D8\
}\
.ace-panda .ace_comment {\
font-style: italic;\
color: #676B79\
}\
.ace-panda .ace_variable {\
color: #FFAAD9\
}\
.ace-panda .ace_variable.ace_parameter {\
font-style: italic;\
color: #C7C7C7\
}\
.ace-panda .ace_entity.ace_name.ace_tag {\
color: #FF2C6D\
}";

var dom = require("../lib/dom");
dom.importCssString(exports.cssText, exports.cssClass);
});
import { defineComponent, h } from "vue";
import { Icon, Icon as IconifyIcon, addCollection, addIcon, listIcons } from "@iconify/vue";
import { ArrowDown, ArrowLeft, ArrowLeftFromLine as MdiMenuOpen, ArrowLeftToLine, ArrowRightFromLine as MdiMenuClose, ArrowRightLeft, ArrowRightToLine, ArrowUp, ArrowUpToLine, Bell, Bold, BookOpenText, Check, ChevronDown, ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight, Circle, CircleAlert, CircleCheckBig, CircleHelp, CircleX, CloudUpload, Copy, CornerDownLeft, Download, Ellipsis, Eraser, Expand, ExternalLink, Eye, EyeOff, FoldHorizontal, Fullscreen, Github, Grid, Grip, GripVertical, Heading1, Heading2, Highlighter, History, ImagePlus, Inbox, Info, InspectionPanel, Italic, Languages, LayoutGrid, Link2, List, ListOrdered, LoaderCircle, LockKeyhole, LogOut, MailCheck, Maximize, Menu, Menu as IconDefault, MessageSquareCode, Minimize, Minimize2, MoonStar, Paintbrush, Palette, PanelLeft, PanelRight, Pin, PinOff, Plus, Redo2, RefreshCw, RemoveFormatting, RotateCw, Search, SearchX, Settings, ShieldQuestion, Shrink, Square, SquareCheckBig, SquareCode, SquareMinus, Strikethrough, Sun, SunMoon, SwatchBook, TextAlignCenter as AlignCenter, TextAlignEnd as AlignRight, TextAlignStart as AlignLeft, TextQuote, Trash2, Underline, Undo2, Unlink2, Upload, UserRoundPen, X } from "lucide-vue-next";
//#region src/create-icon.ts
function createIconifyIcon(icon) {
	return defineComponent({
		name: `Icon-${icon}`,
		setup(props, { attrs }) {
			return () => h(Icon, {
				icon,
				...props,
				...attrs
			});
		}
	});
}
//#endregion
export { AlignCenter, AlignLeft, AlignRight, ArrowDown, ArrowLeft, ArrowLeftToLine, ArrowRightLeft, ArrowRightToLine, ArrowUp, ArrowUpToLine, Bell, Bold, BookOpenText, Check, ChevronDown, ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight, Circle, CircleAlert, CircleCheckBig, CircleHelp, CircleX, CloudUpload, Copy, CornerDownLeft, Download, Ellipsis, Eraser, Expand, ExternalLink, Eye, EyeOff, FoldHorizontal, Fullscreen, Github, Grid, Grip, GripVertical, Heading1, Heading2, Highlighter, History, IconDefault, IconifyIcon, ImagePlus, Inbox, Info, InspectionPanel, Italic, Languages, LayoutGrid, Link2, List, ListOrdered, LoaderCircle, LockKeyhole, LogOut, MailCheck, Maximize, MdiMenuClose, MdiMenuOpen, Menu, MessageSquareCode, Minimize, Minimize2, MoonStar, Paintbrush, Palette, PanelLeft, PanelRight, Pin, PinOff, Plus, Redo2, RefreshCw, RemoveFormatting, RotateCw, Search, SearchX, Settings, ShieldQuestion, Shrink, Square, SquareCheckBig, SquareCode, SquareMinus, Strikethrough, Sun, SunMoon, SwatchBook, TextQuote, Trash2, Underline, Undo2, Unlink2, Upload, UserRoundPen, X, addCollection, addIcon, createIconifyIcon, listIcons };

#include "includes.h"
#include "gtkccallbacksheader.h"
#include "_cgo_export.h"
#include <string.h>
#include <stdio.h>

GtkWindow* Window = NULL;
GtkGLArea* GLarea = NULL;
char* ErrorString = NULL;
gboolean mouseInGLarea = FALSE;

void initialise(int args,char** argv)
{
	ErrorString = (char*)malloc(150);
	gtk_init(&args,&argv);
}

int createWindow(unsigned int width, unsigned int height, const char* title)
{
	Window = createWindowObject();
	configureWindowParameters(Window,width,height,title);
	connectWindowSignals(Window);
	
    createGLArea();

	gtk_widget_show_all(GTK_WIDGET(Window));
	return 1;
}

GtkWindow* createWindowObject()
{
	return (GtkWindow*)gtk_window_new(GTK_WINDOW_TOPLEVEL);
}

void configureWindowParameters(GtkWindow* window,unsigned int width, unsigned int height, const char* title)
{
	if(width != 0 && height != 0)
	{
		gtk_widget_set_size_request(GTK_WIDGET(window),1,1);
		gtk_window_resize(window,width,height);
	}
	if(strcmp(title,"")!=0)
		gtk_window_set_title(window,title);
}

void connectWindowSignals(GtkWindow* window)
{
	g_signal_connect(GTK_WIDGET(window),"delete-event",G_CALLBACK(gtkgo_quit_c),NULL);
	
	g_signal_connect(GTK_WIDGET(window),"key-press-event",G_CALLBACK(gtkgo_gl_area_key_press_c),NULL);
	g_signal_connect(GTK_WIDGET(window),"key-release-event",G_CALLBACK(gtkgo_gl_area_key_release_c),NULL);
	g_signal_connect(GTK_WIDGET(window),"button-press-event",G_CALLBACK(gtkgo_gl_area_button_press_c),NULL);
	g_signal_connect(GTK_WIDGET(window),"button-release-event",G_CALLBACK(gtkgo_gl_area_button_release_c),NULL);
	g_signal_connect(GTK_WIDGET(window),"motion-notify-event",G_CALLBACK(gtkgo_gl_area_motion_notify_c),NULL);
	g_signal_connect(GTK_WIDGET(window),"scroll-event",G_CALLBACK(gtkgo_gl_area_scroll_c),NULL);

	gtk_widget_set_events(GTK_WIDGET(window), GDK_POINTER_MOTION_MASK|GDK_SCROLL_MASK|GDK_ENTER_NOTIFY_MASK);
}

void createGLArea()
{
	GLarea = (GtkGLArea*)gtk_gl_area_new();
	configureGLArea(GLarea);
}

void configureGLArea(GtkGLArea* area)
{
	g_signal_connect(GTK_WIDGET(area),"render",G_CALLBACK(gtkgo_gl_area_render_c),NULL);
	g_signal_connect(GTK_WIDGET(area),"realize",G_CALLBACK(gtkgo_gl_area_realize_c),NULL);
	
	gtk_gl_area_set_has_depth_buffer(area,TRUE);
	gtk_widget_set_events(GTK_WIDGET(area), GDK_POINTER_MOTION_MASK|GDK_SCROLL_MASK|GDK_ENTER_NOTIFY_MASK);

	gdk_threads_add_idle(queue_render_idle,area);

	if(gtk_widget_get_parent(GTK_WIDGET(area)) != NULL)
	{
		gtk_widget_unrealize(GTK_WIDGET(area));
		gtk_widget_realize(GTK_WIDGET(area));
	}
}

void windowSetSize(float width, float height)
{
	gtk_window_resize(Window,width,height);
}

void windowGetSize(float* width, float* height)
{
	if(width == NULL || height == NULL)
	{
		return;
	}

	GtkAllocation* alloc = g_new(GtkAllocation, 1);
    gtk_widget_get_allocation(GTK_WIDGET(Window), alloc);
    if(alloc != NULL)
    {
    	*width = alloc->width;
    	*height = alloc->height;    	
    	g_free(alloc);
    }
    else
    {
    	*width = -1.0;
    	*height = -1.0;
    }
}

void loop()
{
	gtk_main();
}

void windowHideCursor()
{
	gdk_window_set_cursor(gtk_widget_get_window(GTK_WIDGET(Window)),gdk_cursor_new_from_name(gdk_display_get_default(),"none"));
}

void windowDisableCursor()
{
	// gdk_window_set_cursor(gtk_widget_get_window(GTK_WIDGET(Window)),gdk_cursor_new_from_name(gdk_display_get_default(),"none"));
	if(gdk_seat_grab(
				gdk_display_get_default_seat(gdk_display_get_default()),
				gtk_widget_get_window(GTK_WIDGET(Window)),
				GDK_SEAT_CAPABILITY_ALL_POINTING,
				FALSE,
				gdk_cursor_new_from_name(gdk_display_get_default(),"none"),
				NULL,
				NULL,
				NULL
	) != GDK_GRAB_SUCCESS)
	{
		g_print("Error disabling cursor\n");
	}
}

void windowShowCursor()
{
	gdk_seat_ungrab(gdk_display_get_default_seat(gdk_display_get_default()));
	gdk_window_set_cursor(gtk_widget_get_window(GTK_WIDGET(Window)),gdk_cursor_new_for_display(gdk_display_get_default(),GDK_ARROW));
}

int windowCursorShown()
{
	GdkCursor* cursor = gdk_window_get_cursor(gtk_widget_get_window(GTK_WIDGET(Window)));
	return cursor != NULL && gdk_cursor_get_cursor_type(cursor) != GDK_BLANK_CURSOR;
}
int windowCursorHidden()
{
	GdkCursor* cursor = gdk_window_get_cursor(gtk_widget_get_window(GTK_WIDGET(Window)));
	return cursor == NULL || gdk_cursor_get_cursor_type(cursor) == GDK_BLANK_CURSOR;
}
int windowCursorDisabled()
{
	return windowCursorHidden();
}

void addGLAreaToWindow()
{
	addGLAreaToContainer(GTK_CONTAINER(Window));
}

void addGLAreaToContainer(GtkContainer* container)
{
	gtk_container_add(container,GTK_WIDGET(GLarea));
	gtk_widget_show(GTK_WIDGET(GLarea));
}

GtkContainer* boxToContainer(GtkBox* box)
{
	return GTK_CONTAINER(box);
}

GtkContainer* windowToContainer(GtkWindow* window)
{
	return GTK_CONTAINER(window);
}

GtkWidget* boxToWidget(GtkBox* box)
{
	return GTK_WIDGET(box);
}

GtkBox* widgetToBox(GtkWidget* widget)
{
	return GTK_BOX(widget);
}

GtkWidget* glareaToWidget(GtkGLArea* area)
{
	return GTK_WIDGET(area);
}

GtkWidget* buttonToWidget(GtkButton* button)
{
    return GTK_WIDGET(button);
}

GtkContainer* buttonToContainer(GtkButton* button)
{
    return GTK_CONTAINER(button);
}

GtkButton* widgetToButton(GtkWidget* widget)
{
    return GTK_BUTTON(widget);
}

GtkWidget* gobjectToWidget(GObject* object)
{
	return GTK_WIDGET(object);
}

GObject* widgetToGObject(GtkWidget* widget)
{
	return G_OBJECT(widget);
}

GtkWindow* widgetToWindow(GtkWidget* widget)
{
	return GTK_WINDOW(widget);
}

GtkWidget* gpointerToWidget(gpointer data)
{
	return GTK_WIDGET(data);
}

GtkContainer* widgetToContainer(GtkWidget* widget)
{
	return GTK_CONTAINER(widget);
}

GtkGrid* widgetToGrid(GtkWidget* widget)
{
	return GTK_GRID(widget);
}

GtkWidget* windowToWidget(GtkWindow* window)
{
	return GTK_WIDGET(window);
}

GtkGLArea* gobjectToGLArea(GObject* object)
{
	return GTK_GL_AREA(object);
}

GtkListBox* widgetToListBox(GtkWidget* widget)
{
	return GTK_LIST_BOX(widget);
}

GtkLabel* widgetToLabel(GtkWidget* widget)
{
	return GTK_LABEL(widget);
}

GtkWidget* labelToWidget(GtkLabel* label) 
{
	return GTK_WIDGET(label); 
}

GtkListBox* gobjectToListBox(GObject* object)
{
	return GTK_LIST_BOX(object);
}

GtkMenuItem* gobjectToMenuItem(GObject* object)
{
	return GTK_MENU_ITEM(object);
}

GtkWidget* menuItemToWidget(GtkMenuItem* menuItem)
{
	return GTK_WIDGET(menuItem);
}

GtkGLArea* widgetToGLArea(GtkWidget* widget)
{
	return GTK_GL_AREA(widget);
}

GtkWidget* listBoxToWidget(GtkListBox* listBox)
{
	return GTK_WIDGET(listBox);
}

GtkContainer* listBoxToContainer(GtkListBox* listBox)
{
	return GTK_CONTAINER(listBox);
}

GtkWidget* listBoxRowToWidget(GtkListBoxRow* listBoxRow)
{
	return GTK_WIDGET(listBoxRow);
}

GtkContainer* listBoxRowToContainer(GtkListBoxRow* listBoxRow)
{
	return GTK_CONTAINER(listBoxRow);
}

GtkWidget* toolButtonToWidget(GtkToolButton* toolButton)
{
	return GTK_WIDGET(toolButton);
}

GtkToolButton* gobjectToToolButton(GObject* object)
{
	return GTK_TOOL_BUTTON(object);
}

GtkBox* gobjectToBox(GObject* object) 
{
	return GTK_BOX(object);
}

GtkFileChooserDialog* widgetToFileChooserDialog(GtkWidget* widget)
{
	return GTK_FILE_CHOOSER_DIALOG(widget);
}

GtkDialog* fileChooserDialogToDialog(GtkFileChooserDialog* fcdialog)
{
	return GTK_DIALOG(fcdialog);
}

GtkWidget* fileChooserDialogToWidget(GtkFileChooserDialog* fcdialog)
{
	return GTK_WIDGET(fcdialog);
}

GtkWidget* dialogToWidget(GtkDialog* dialog)
{
	return GTK_WIDGET(dialog);
}

GtkFileChooser* fileChooserDialogToFileChooser(GtkFileChooserDialog* fcdialog)
{
	return GTK_FILE_CHOOSER(fcdialog);
}

GObject* labelToGObject(GtkLabel* label)
{
	return G_OBJECT(label);
}

const gchar* gpointerToGChar(const gpointer data)
{
	return (const gchar*)data;
}

const gpointer gcharToGPointer(const gchar* data)
{
	return (const gpointer)data;
}

GtkWidget* imageToWidget(GtkImage* img)
{
	return GTK_WIDGET(img);
}

GtkImage* widgetToImage(GtkWidget* widget)
{
	return GTK_IMAGE(widget);
}

void widgetGetSize(GtkWidget* widget,gint* width, gint* height)
{
	GtkAllocation* alloc = g_new(GtkAllocation,1);
	gtk_widget_get_allocation(widget,alloc);
	*width = alloc->width;
	*height = alloc->height;
	g_free(alloc);
}

void signalConnectButton(GtkButton* button,const char* signal, int id)
{
    ButtonSignalUserData* bsud = (ButtonSignalUserData*)malloc(sizeof(ButtonSignalUserData));
    bsud->id = id;
    bsud->signal = (char*)malloc(strlen(signal));
	strcpy(bsud->signal,signal);
    g_signal_connect(GTK_WIDGET(button),signal,G_CALLBACK(gtkgo_button_signal_c),bsud);
}

void sizeAllocateSignalConnectWidget(GtkWidget* widget,const char* signal,const char* name)
{
	WidgetSignalUserData* wsud = (WidgetSignalUserData*)malloc(sizeof(WidgetSignalUserData));
	const gchar* namec = gtk_widget_get_name(widget);
	wsud->name = (char*)malloc(strlen(namec));
	strcpy(wsud->name,namec);
	wsud->signal = (char*)malloc(strlen(signal));
	strcpy(wsud->signal,signal);
	g_signal_connect(widget,signal,G_CALLBACK(gtkgo_widget_size_allocate_signal_c),wsud);
}

void signalConnectMenuItem(GtkMenuItem* menuItem,const char* signal,const char* name)
{
	WidgetSignalUserData* wsud = (WidgetSignalUserData*)malloc(sizeof(WidgetSignalUserData));
	const gchar* namec = gtk_widget_get_name(GTK_WIDGET(menuItem));
	wsud->name = (char*)malloc(strlen(namec));
	strcpy(wsud->name,namec);
	wsud->signal = (char*)malloc(strlen(signal));
	strcpy(wsud->signal,signal);
	g_signal_connect(GTK_WIDGET(menuItem),signal,G_CALLBACK(gtkgo_menu_item_signal_c),wsud);

}

void eventSignalConnectWidget(GtkWidget* widget,const char* signal, const char* name)
{
	WidgetSignalUserData* wsud = (WidgetSignalUserData*)malloc(sizeof(WidgetSignalUserData));
	const gchar* namec = gtk_widget_get_name(widget);
	wsud->name = (char*)malloc(strlen(namec));
	strcpy(wsud->name,namec);
	wsud->signal = (char*)malloc(strlen(signal));
	strcpy(wsud->signal,signal);
	g_signal_connect(widget,signal,G_CALLBACK(gtkgo_widget_event_signal_c),wsud);
}

void rowSelectedSignalConnectListBox(GtkListBox* listBox, const char* signal, const char* name)
{
	WidgetSignalUserData* wsud = (WidgetSignalUserData*)malloc(sizeof(WidgetSignalUserData));
	const gchar* namec = gtk_widget_get_name(GTK_WIDGET(listBox));
	wsud->name = (char*)malloc(strlen(namec));
	strcpy(wsud->name,namec);
	wsud->signal = (char*)malloc(strlen(signal));
	strcpy(wsud->signal,signal);
	g_signal_connect(GTK_WIDGET(listBox),signal,G_CALLBACK(gtkgo_list_box_row_selected_signal_c),wsud);
}

void signalConnectToolButton(GtkToolButton* toolButton, const char* name)
{
	char* namec = (char*)malloc(strlen(name));
	strcpy(namec,name);
	g_signal_connect(GTK_WIDGET(toolButton),"clicked",G_CALLBACK(gtkgo_tool_button_signal_c),namec);
}

GtkWidget* gohome_file_chooser_dialog_new(const gchar *title,GtkWindow *parent,GtkFileChooserAction action)
{
	GtkWidget* widget;
	switch(action)
	{
		default:
			widget = gtk_file_chooser_dialog_new(title,parent,action,
				"gtk-cancel",GTK_RESPONSE_CANCEL,
				"gtk-open",GTK_RESPONSE_ACCEPT,
				NULL
			);
			break;
	}
	return widget;
}




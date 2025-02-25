package com.example;

import org.telegram.telegrambots.meta.api.methods.send.SendMessage;
import org.telegram.telegrambots.meta.api.methods.updatingmessages.EditMessageText;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.InlineKeyboardMarkup;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.buttons.InlineKeyboardButton;

import java.util.Collections;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

public class AdminCommands{
	private Words words = new Words();
	private SendMessage message = new SendMessage();

	public SendMessage dataBase(String messageText, long chatId){
		message.setChatId(String.valueOf(chatId));
		message.setReplyMarkup(null);
		String[] parts = messageText.split(" ", 5);
		if(messageText.contains("/add")){
			String task_id = parts[2];
			String wordName = parts[3];
			String wordExplanation = parts[4];
			message.setText(words.addWord(wordName, wordExplanation, task_id));
		}
		else if(messageText.contains("/remove")){
			String task_id = parts[2];
			String wordName = parts[3];
			message.setText(words.removeWord(wordName, task_id));

		}
		else if(messageText.contains("/showall")){
			UserState userState = UserStateManager.getUserState(chatId);
			String task_id = parts[2];
			userState.allWords = words.showAllBase(task_id);
			userState.isSetting = true;
			userState.currentPage = 0;
			message.setText(userState.allWords.get(userState.currentPage));
			
			InlineKeyboardMarkup showKeyboard = new InlineKeyboardMarkup();
			InlineKeyboardButton back = new InlineKeyboardButton(); back.setText("<");
			InlineKeyboardButton next = new InlineKeyboardButton(); next.setText(">");	
    	InlineKeyboardButton allPages = new InlineKeyboardButton(); allPages.setText("[" + (userState.currentPage+1) + "/" + (userState.allWords.size()) + "]");
			back.setCallbackData("back"); next.setCallbackData("next"); allPages.setCallbackData("toMain");
			
			List<InlineKeyboardButton> row = Arrays.asList(back, next);
			showKeyboard.setKeyboard(Arrays.asList(
					row,
					Collections.singletonList(allPages)
					));
			message.setReplyMarkup(showKeyboard);
		}
		else if(messageText.contains("/adm")){
			UserStateManager.getUserState(chatId).isSetting = false; 
			message.setText("Текущие админ команды: " +
					"\n/add [task number] [word explanation + .] - Добавляет слово к заданию [task number] с пояснением [word explanation] где точка - конец пояснения" +
					"\n/remove [task number] [word] - Удаляет слово [word] из задания [task number]" +
					"\n/showall [task number] - Просмотр всех слов к заданию [task number]"
					);

		}
		else{
			String text = "Команды: " + "\"" + parts[1] + "\"" +  " не существует";
			message.setText(text);
		}
		return message;
	}
	public EditMessageText showBack(long callbackChatId, long callbackMessageId){
    UserState userState = UserStateManager.getUserState(callbackChatId);
		userState.currentPage = (userState.currentPage > 0) ? userState.currentPage-1 : 0;
    EditMessageText editMessage = createEditMessage(callbackChatId, callbackMessageId);
		editMessage.setText(userState.allWords.get(userState.currentPage));
		return editMessage;
	}
	public EditMessageText showNext(long callbackChatId, long callbackMessageId){
		UserState userState = UserStateManager.getUserState(callbackChatId);
    userState.currentPage = (userState.currentPage < userState.allWords.size()) ? userState.currentPage+1 : userState.allWords.size()-1;
    EditMessageText editMessage = createEditMessage(callbackChatId, callbackMessageId);
    editMessage.setText(userState.allWords.get(userState.currentPage));
  	return editMessage;
	}
	private EditMessageText createEditMessage (long callbackChatId, long callbackMessageId){
		UserState userState = UserStateManager.getUserState(callbackChatId);
		InlineKeyboardMarkup showKeyboard = new InlineKeyboardMarkup();
    InlineKeyboardButton back = new InlineKeyboardButton(); back.setText("<");
    InlineKeyboardButton next = new InlineKeyboardButton(); next.setText(">");
    InlineKeyboardButton allPages = new InlineKeyboardButton(); allPages.setText("[" + (userState.currentPage+1) + "/" + (userState.allWords.size()) + "]");
   	back.setCallbackData("back"); next.setCallbackData("next"); allPages.setCallbackData("toMain");
		List<InlineKeyboardButton> row = Arrays.asList(back, next);
    showKeyboard.setKeyboard(Arrays.asList(
					row,
					Collections.singletonList(allPages)
					));
		
		EditMessageText editMessage = new EditMessageText();
    editMessage.setChatId(String.valueOf(callbackChatId));
    editMessage.setMessageId((int) callbackMessageId);
    editMessage.setReplyMarkup(showKeyboard);
    editMessage.setText("Error 404");
	
  	return editMessage;
	}
}

package com.example;

import java.sql.PreparedStatement;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.sql.Connection;
import java.sql.ResultSet;

import org.telegram.telegrambots.meta.api.objects.Update;
import org.telegram.telegrambots.meta.api.objects.Message;
import org.telegram.telegrambots.bots.TelegramLongPollingBot;
import org.telegram.telegrambots.meta.api.objects.CallbackQuery;
import org.telegram.telegrambots.meta.api.methods.send.SendMessage;
import org.telegram.telegrambots.meta.exceptions.TelegramApiException;

import org.telegram.telegrambots.meta.api.objects.replykeyboard.ReplyKeyboardMarkup;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.buttons.KeyboardRow;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.ReplyKeyboardRemove;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.buttons.KeyboardButton;

import org.telegram.telegrambots.meta.api.objects.replykeyboard.InlineKeyboardMarkup;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.buttons.InlineKeyboardButton;


import java.util.List;
import java.util.Arrays;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Collections;


public class ReplyKeyboard{
	private MainKeyboard main = new MainKeyboard();
	private CheckKeyboard check = new CheckKeyboard();

	private Statistic statistic = new Statistic();
	private SendMessage messageMenu = new SendMessage();
	private MainFunctional functional = new MainFunctional();
	
	public SendMessage createMenu(Message message, long chatId, long userId){
		String messageText = "";
		String userName = "";

		if(message != null){
			messageText = message.getText();
			userName = message.getFrom().getFirstName();
		}

		UserState userState = UserStateManager.getUserState(userId);
		if(message != null && message.getFrom().getLastName() != null) {userName += message.getFrom().getLastName();}

		messageMenu.setChatId(String.valueOf(chatId));
		messageMenu.setReplyMarkup(null);

		if(messageText.equals("Выбрать задание") && !userState.isChoosing && !userState.isChecking){
			messageMenu = check.createMenu(chatId, userId);
			userState.isChoosing = true;
		}
		else if(messageText.equals("Проверить знания") && !userState.isChoosing){
			messageMenu = functional.makeTask(userId);
			if(messageMenu.getText() != "Такого задания ещё нет в RusBooster") {
				userState.isChecking = true;
			}
			else{
				messageMenu = main.createMenu(chatId);
				messageMenu.setText("Такого задания ещё нет в RusBooster");
			}
		}
		else if(messageText.equals("Статистика") && !userState.isChecking && !userState.isChoosing){
			messageMenu.setText(statistic.getStatistic(userName, userId));
		}

		else{
			messageMenu = main.createMenu(chatId);
			if(userState.isChoosing){
				userState.isChoosing = false;
				messageMenu.setText(statistic.chooseExercise(userId, messageText));
				if(messageMenu.getText().equals("Перенапровляемся...")){
					messageMenu = main.createMenu(chatId);
				}
			}
			
			if(userState.isChecking){
				userState.isChecking = false;
				messageMenu = functional.explanationTask(chatId, userId, messageText);
			}
		}
		return messageMenu;
	}
}

class MainKeyboard{
	public SendMessage createMenu(long chatId){
		SendMessage message = new SendMessage();
		message.setChatId(String.valueOf(chatId));
		message.setText("Выберите опцию: ");

		ReplyKeyboardMarkup keyboardMarkup = new ReplyKeyboardMarkup();
		keyboardMarkup.setResizeKeyboard(true);

		List<KeyboardRow> keyboard = new ArrayList<>();
		KeyboardRow row = new KeyboardRow();
		row.add(new KeyboardButton("Выбрать задание"));
		row.add(new KeyboardButton("Проверить знания"));
		row.add(new KeyboardButton("Статистика"));
		keyboard.add(row);

		keyboardMarkup.setKeyboard(keyboard);
		message.setReplyMarkup(keyboardMarkup);

		return message;
	}
}

class CheckKeyboard{
	private String url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/statistic.db";
	private String sql;

	public SendMessage createMenu(long chatId, long userId){
		SendMessage message = new SendMessage();
		message.setChatId(String.valueOf(chatId));
		message.setText("Напишите номер упражнения от 1 до 26: ");

		InlineKeyboardMarkup keyboardMarkup = new InlineKeyboardMarkup();
		InlineKeyboardButton baddestTask = new InlineKeyboardButton();
		InlineKeyboardButton cancelChoose = new InlineKeyboardButton();

		cancelChoose.setText("В главное меню");
		cancelChoose.setCallbackData("cancelChoose");
		baddestTask.setCallbackData("baddestTask");

		try(Connection conn = DriverManager.getConnection(url)){
			sql = "SELECT baddest_task FROM statistics WHERE user_id = ?";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setLong(1, userId);
			ResultSet result = pstmt.executeQuery();
			if(result.next()){
				if(result.getInt("baddest_task") != 0){
					baddestTask.setText("Наихудшая успеваимость: №" + result.getInt("baddest_task"));
					keyboardMarkup.setKeyboard(Arrays.asList(
						Collections.singletonList(baddestTask),
						Collections.singletonList(cancelChoose)
					));
				}
			}
			else{
				keyboardMarkup.setKeyboard(Collections.singletonList(Collections.singletonList(cancelChoose)));
			}
		}
		catch(SQLException e){
			e.printStackTrace();
		}
		message.setReplyMarkup(keyboardMarkup);
		return message;
	}
}

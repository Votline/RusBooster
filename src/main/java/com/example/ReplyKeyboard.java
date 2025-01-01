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

import java.util.List;
import java.util.ArrayList;
import java.io.IOException;

public class ReplyKeyboard{
	private MainKeyboard main = new MainKeyboard();
	private CheckKeyboard check = new CheckKeyboard();

	private Settings settings = new Settings();
	private Statistic statistic = new Statistic();
	private SendMessage messageMenu = new SendMessage();
	private MainFunctional functional = new MainFunctional();
	
	public void createMenu(TelegramLongPollingBot bot, Message message){
		String messageText = message.getText();

		long chatId = message.getChatId();
		long userId = message.getFrom().getId();
		String userName = message.getFrom().getFirstName();
		
		UserState userState = UserStateManager.getUserState(userId);
		if(message.getFrom().getLastName() != null) userName += message.getFrom().getLastName();

		this.messageMenu.setChatId(String.valueOf(chatId));
		this.messageMenu.setReplyMarkup(null);

		if(messageText.equals("Выбрать задание") && !userState.isChoosing && !userState.isChecking){
			this.messageMenu = check.createMenu(chatId, userId);
			userState.isChoosing = true;
		}
		else if(messageText.equals("Проверить знания") && !userState.isChoosing){
			String textForMenu = functional.makeTask(userId);
			if(textForMenu != "Такого задания ещё нет в RusBooster") {
				userState.isChecking = true;
				ReplyKeyboardRemove keyboardRemove = new ReplyKeyboardRemove(true);
				this.messageMenu.setReplyMarkup(keyboardRemove);
			}
			this.messageMenu.setText(textForMenu);
		}
		else if(messageText.equals("Статистика") && !userState.isChecking && !userState.isChoosing){
			this.messageMenu.setText(statistic.getStatistic(userName, userId));
		}
		else{
			this.messageMenu = main.createMenu(chatId);
			try{
				if(userState.isChoosing){
					userState.isChoosing = false;
					settings.chooseExercise(userId, messageText);
				}
			}
			catch(NumberFormatException e){
				e.printStackTrace();
			}
			if(userState.isChecking){
				userState.isChecking = false;
				this.messageMenu = functional.explanationTask(chatId, userId, messageText);
			}
		}

		try{bot.execute(this.messageMenu); }
		catch(TelegramApiException e){e.printStackTrace(); }
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

		ReplyKeyboardMarkup keyboardMarkup = new ReplyKeyboardMarkup();
		keyboardMarkup.setResizeKeyboard(true);
		List<KeyboardRow> keyboard = new ArrayList<>();
		KeyboardRow row2 = new KeyboardRow();

		try(Connection conn = DriverManager.getConnection(url)){
			sql = "SELECT current_task FROM statistics WHERE user_id = ?";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setLong(1, userId);
			ResultSet result = pstmt.executeQuery();
			if(result.next()){
				KeyboardRow row1 = new KeyboardRow();
				row1.add(new KeyboardButton("Текущее задание: " + result.getString("current_task") ));
				keyboard.add(row1);
			}
		}
		catch(SQLException e){
			e.printStackTrace();
		}
		row2.add(new KeyboardButton("Назад"));
		keyboard.add(row2);

		keyboardMarkup.setKeyboard(keyboard);
		message.setReplyMarkup(keyboardMarkup);

		return message;
	}
}

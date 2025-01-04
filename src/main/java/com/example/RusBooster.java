import java.util.List;
import java.nio.file.*;

import org.telegram.telegrambots.meta.api.objects.User;
import org.telegram.telegrambots.meta.api.objects.Update;
import org.telegram.telegrambots.meta.api.objects.Message;
import org.telegram.telegrambots.bots.TelegramLongPollingBot;
import org.telegram.telegrambots.meta.api.methods.send.SendMessage;
import org.telegram.telegrambots.meta.exceptions.TelegramApiException;


public class RusBooster extends TelegramLongPollingBot{
	private static final String BOT_TOKEN = loadToken();
	private static final String BOT_NAME = "RusBooster";

	private long chatId;
	private long userId;
	private String messageText;

	ReplyKeyboard botMenu = new ReplyKeyboard();
	AdminCommands adminCommands = new AdminCommands();
	MainFunctional functional = new MainFunctional();
	
	@Override
	public String getBotUsername(){
		return BOT_NAME;
	}

	@Override
	public String getBotToken(){
		return BOT_TOKEN;
	}

	@Override
	public void onUpdateReceived(Update update){

		if(update.hasMessage() && update.getMessage().hasText()){
			messageText = update.getMessage().getText();
			userId = update.getMessage().getFrom().getId();
			chatId = update.getMessage().getChatId();
		}
		
		if(update.hasCallbackQuery() && update.getCallbackQuery().getData() != null){	
			UserState userState = UserStateManager.getUserState(userId);

			String callbackData = update.getCallbackQuery().getData();
			long callbackChatId = update.getCallbackQuery().getMessage().getChatId();
			long callbackUserId = update.getCallbackQuery().getFrom().getId();
			
			if("showExplanations".equals(callbackData)){
				try{this.execute(functional.showExplanations(callbackChatId));}
				catch(TelegramApiException e){e.printStackTrace();}
				return;
			}
			else if("cancelTask".equals(callbackData)){
				userState.isChecking = false;
			}
			try{this.execute(botMenu.createMenu(update.getMessage(), callbackChatId, callbackUserId));}	
			catch(TelegramApiException e){e.printStackTrace();}

			return;
		}

		if(messageText.contains("/adm") && userId == 5459965917L){
			try{this.execute(adminCommands.dataBase(messageText, chatId));}
			catch(TelegramApiException e){e.printStackTrace();}
		}
		else if(messageText != null){
			try{this.execute(botMenu.createMenu(update.getMessage(), chatId, userId));}
			catch(TelegramApiException e){e.printStackTrace();}
		}
	}
	private static String loadToken() {
        try {
            return Files.readString(Path.of("token.txt")).trim();
        } catch (Exception e) {
            throw new RuntimeException("Ошибка загрузки токена из файла", e);
        }
    }
}

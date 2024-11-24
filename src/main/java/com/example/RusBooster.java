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

	ReplyKeyboard botMenu = new ReplyKeyboard();
	AdminCommands adminCommands = new AdminCommands();

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
		String messageText = update.getMessage().getText();
		long userId = update.getMessage().getFrom().getId();
		long chatId = update.getMessage().getChatId();


		if(messageText.contains("/adm") && userId == 5459965917L){
			try{
				this.execute(adminCommands.dataBase(messageText, chatId));
			}
			catch(TelegramApiException e){
				e.printStackTrace();
			}
		}
		else if(messageText != null){
			botMenu.createMenu(this, update.getMessage());
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

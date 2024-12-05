package;

import libcron.Scheduler;
import Miniaudio;
import logging.Logger;

class Bell {
	public static var instance: Bell;
	var log: Logger;
	var cron: libcron.Scheduler = new libcron.Scheduler();
	var dataPointer: cpp.Star<cpp.UInt8>;
	var decoder: MaDecoder;
	var engine: MaEngine;
	var sound: MaSound;

	public function new() {
		this.log = new Logger(Bell);
		this.cron = new Scheduler();

		// Decode the bell mp3
		final bell = haxe.Resource.getBytes("bell_mp3");
		this.dataPointer = cpp.Pointer.arrayElem(bell.getData(), 0).ptr;
		this.decoder = Miniaudio.MaDecoder.create();
		final config = Miniaudio.ma_decoder_config_init_default();
		final result = Miniaudio.ma_decoder_init_memory(cast dataPointer, bell.length, config, decoder);
		if (result != MaResult.MA_SUCCESS) {
			trace("Failed to decode audio");
			return;
		}

		// Create a miniaudio engine
		this.engine = MaEngine.create();
		final result = Miniaudio.ma_engine_init(null, engine);
		if (result != MaResult.MA_SUCCESS) {
			trace("Failed to initialize engine");
			return;
		}

		// Create the sound and play it with the engine
		this.sound = MaSound.create();
		final result = Miniaudio.ma_sound_init_from_data_source(engine, cast decoder, 0, null, sound);
		if (result != MaResult.MA_SUCCESS) {
			trace("Failed to initialize sound");
			return;
		}

		// Tick cron
		var timer = new haxe.Timer(1000);
		timer.run = () -> {
			instance.cron.tick();
		};
	}

	public function ring(?job: Int) {
		log.info("Ring!");
		Miniaudio.ma_sound_start(sound);
	}

	static function instanceRing() {
		instance.ring();
	}

	public static function schedule(cron: String) {
		trace(cron);
		trace(instance.cron.addSchedule(cpp.StdString.ofString("Bell"), cpp.StdString.ofString(cron), cpp.Callable.fromStaticFunction(instanceRing)));
	}
}
